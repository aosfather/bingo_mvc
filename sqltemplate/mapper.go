package sqltemplate

import (
	"bytes"
	"container/list"
	"fmt"
	utils "github.com/aosfather/bingo_utils"
	"github.com/aosfather/bingo_utils/files"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"sync"
	"text/template"
)

/*
  mapper dao
  基于yaml文件
  方法：find count exist query update insert delete
*/
type methodType byte

const (
	mt_find   methodType = 1
	mt_count  methodType = 2
	mt_exist  methodType = 3
	mt_query  methodType = 4
	mt_update methodType = 10
	mt_insert methodType = 11
	mt_delete methodType = 12
)

func (this *methodType) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	text = strings.ToLower(text)
	if text == "find" {
		*this = mt_find
	} else if text == "count" {
		*this = mt_count
	} else if text == "exist" {
		*this = mt_exist
	} else if text == "query" {
		*this = mt_query
	} else if text == "update" {
		*this = mt_update
	} else if text == "insert" {
		*this = mt_insert
	} else if text == "delete" {
		*this = mt_delete
	} else {
		*this = 0
		return fmt.Errorf("value is wrong! [ %s ]", text)
	}
	return nil
}

type MapperFile struct {
	Name  string `yaml:"namespace"`
	Nodes []MapperNode
}

type MapperNode struct {
	Code     string
	Type     methodType
	Template string
}

func (this *MapperNode) ToFunction(namespace string) mapperfunction {
	f := mapperfunction{Code: this.Code, Type: this.Type}
	f.Init(namespace, this.Template)
	return f
}

type mapperfunction struct {
	Code string
	Type methodType
	t    *template.Template
	args *list.List
	lock sync.Mutex
}

func (this *mapperfunction) Init(namespace string, temp string) {
	this.args = list.New()
	this.t = template.New(fmt.Sprintf("%s::%s", namespace, this.Code))
	this.t.Funcs(template.FuncMap{"sql": this.call})
	this.t.Parse(temp)
}
func (this *mapperfunction) call(str reflect.Value) interface{} {
	this.args.PushBack(str.Interface())
	return "?"

}

func (this *mapperfunction) CreateSql(input utils.Object) (string, []interface{}) {
	if this.t != nil {
		this.lock.Lock()
		//重新初始化列表
		this.args.Init()
		buffer := new(bytes.Buffer)
		err := this.t.Execute(buffer, input)
		if err != nil {
			buffer.WriteString(err.Error())
		}

		//将列表中的数据拷贝到参数列表中
		var args []interface{}
		for e := this.args.Front(); e != nil; e = e.Next() {
			args = append(args, e.Value)
		}

		defer this.lock.Unlock()
		return buffer.String(), args

	}
	return fmt.Sprintf("error:%s not exist!", this.Code), nil
}

type SqltemplateCollect map[string]mapperfunction
type SqltemplateManager struct {
	templateCollects map[string]SqltemplateCollect
}

func (this *SqltemplateManager) AddCollectFromFile(f string) {
	if this.templateCollects == nil {
		this.templateCollects = make(map[string]SqltemplateCollect)
	}
	mapperfile := MapperFile{}
	if files.IsFileExist(f) {
		yamlFile, err := ioutil.ReadFile(f)
		if err != nil {
			log.Println(err.Error())
			return
		}
		yaml.Unmarshal(yamlFile, &mapperfile)

		//构建mappfunction
		namespace := mapperfile.Name
		collect := make(SqltemplateCollect)
		for _, node := range mapperfile.Nodes {
			collect[node.Code] = node.ToFunction(namespace)
		}
		this.templateCollects[namespace] = collect

	}

}

func (this *SqltemplateManager) BuildDao(ds *DataSource, namespace string) *MapperDao {
	if v, ok := this.templateCollects[namespace]; ok {
		return &MapperDao{ds, v, nil}
	}
	return nil
}

type MapperDao struct {
	ds        *DataSource
	templates SqltemplateCollect
	current   *Connection
}

func (this *MapperDao) BeginTransaction() {
	this.current = this.ds.GetConnection()
	this.current.Begin()
}

func (this *MapperDao) FinishTransaction() {
	if this.current != nil {
		this.current.Commit()
		this.current.Close()
		this.current = nil
	}
}

func (this *MapperDao) GetDao() *BaseDao {
	return this.ds.GetDao()
}

type exec func(conn *Connection, sqlstring string, args []interface{})

func (this *MapperDao) execute(obj utils.Object, id string, mt methodType, e exec) {
	session := this.current
	if session == nil {
		session = this.ds.GetConnection()
		if mt == mt_delete || mt == mt_insert || mt == mt_update {
			session.Begin()
		}
		defer session.Close()
	}

	//如果给的id为空则说明不走模板引擎
	if id == "" {
		if e != nil {
			e(session, "", nil)
		}

		return
	}
	//根据id查找模板
	function := this.templates[id]
	if function.Type == mt {
		sqlstr, args := function.CreateSql(obj)
		if e != nil {
			e(session, sqlstr, args)
		}
	}

}
func (this *MapperDao) Find(obj utils.Object, id string) bool {
	var result bool
	fun := func(conn *Connection, sqlstring string, args []interface{}) {
		log.Println(sqlstring)
		log.Println(args)
		result = conn.Query(obj, sqlstring, args...)
	}
	this.execute(obj, id, mt_find, fun)
	return result
}

func (this *MapperDao) executeCommand(obj utils.Object, command methodType, id string) (int64, int64, error) {
	var dbId, affect int64
	var err error
	fun := func(conn *Connection, sqlstring string, args []interface{}) {
		if sqlstring == "" {
			if command == mt_insert {
				dbId, affect, err = this.current.Insert(obj)
			}

		} else {
			dbId, affect, err = conn.ExeSql(sqlstring, args...)
		}

		if err != nil {
			conn.Rollback()
		} else {
			if this.current == nil {
				conn.Commit()
			}
		}
	}

	this.execute(obj, id, command, fun)
	return dbId, affect, err
}

func (this *MapperDao) Insert(obj utils.Object, id string) (int64, error) {
	dbId, _, err := this.executeCommand(obj, mt_insert, id)
	return dbId, err
}

//一般的insert只有一条，自动选择一条，如果不存在insert，将使用template来完成插入
func (this *MapperDao) InsertAuto(obj utils.Object) (int64, error) {
	//查找第一条insert模板
	for _, v := range this.templates {
		if v.Type == mt_insert {
			return this.Insert(obj, v.Code)
		}
	}
	//如果没有insert语句，使用自动template插入
	dbId, _, err := this.executeCommand(obj, mt_insert, "")
	return dbId, err

}

func (this *MapperDao) Update(obj utils.Object, id string) (int64, error) {
	_, affect, err := this.executeCommand(obj, mt_update, id)
	return affect, err
}

func (this *MapperDao) Delete(obj utils.Object, id string) (int64, error) {
	_, affect, err := this.executeCommand(obj, mt_delete, id)
	return affect, err
}

func (this *MapperDao) Query(obj utils.Object, page Page, id string) []interface{} {
	var result []interface{}
	fun := func(conn *Connection, sqlstring string, args []interface{}) {
		result = conn.QueryByPage(obj, page, sqlstring, args...)
	}
	this.execute(obj, id, mt_query, fun)
	return result

}

func (this *MapperDao) Count(obj utils.Object, id string) int64 {
	var result int64
	fun := func(conn *Connection, sqlstring string, args []interface{}) {
		b := conn.Query(&result, sqlstring, args...)
		if !b {
			result = 0
		}
	}
	this.execute(obj, id, mt_count, fun)
	return result
}

func (this *MapperDao) Exist(obj utils.Object, id string) bool {
	var result bool
	fun := func(conn *Connection, sqlstring string, args []interface{}) {
		lowcase := strings.ToLower(sqlstring)
		if strings.Index(lowcase, "select") <= 0 {
			sqlstring = fmt.Sprintf("select 1 %s", sqlstring)
		}
		var ds int64
		result = conn.Query(&ds, sqlstring, args)
	}
	this.execute(obj, id, mt_exist, fun)
	return result

}
