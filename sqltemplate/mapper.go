package sqltemplate

import (
	"bytes"
	"fmt"
	utils "github.com/aosfather/bingo_utils"
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

type MapperFile struct {
	Name  string
	Nodes []MapperNode
}

type MapperNode struct {
	Code     string
	Type     methodType
	Template string
}

type mapperfunction struct {
	Code string
	Type methodType
	t    *template.Template
}

func (this *mapperfunction) CreateSql(input utils.Object) string {
	if this.t != nil {
		buffer := new(bytes.Buffer)
		err := this.t.Execute(buffer, input)
		if err != nil {
			buffer.WriteString(err.Error())
		}
		return buffer.String()
	}
	return fmt.Sprintf("error:%s not exist!", this.Code)
}

type MapperDao struct {
	ds        *DataSource
	templates map[string]mapperfunction
	current   *Connection
}

func (this *MapperDao) loadFromfile() {

}

func (this *MapperDao) BeginTransaction() {
	this.current = this.ds.GetConnection()
	this.current.Begin()
}

func (this *MapperDao) GetDao() *BaseDao {
	return this.ds.GetDao()
}

type exec func(conn *Connection, sqlstring string)

func (this *MapperDao) execute(obj utils.Object, id string, mt methodType, e exec) {
	session := this.current
	if session == nil {
		session = this.ds.GetConnection()
		defer session.Close()
	}
	function := this.templates[id]
	if function.Type == mt {
		sqlstr := function.CreateSql(obj)
		if e != nil {
			e(session, sqlstr)
		}
	}

}
func (this *MapperDao) Find(obj utils.Object, id string) bool {
	var result bool
	fun := func(conn *Connection, sqlstring string) {
		result = conn.Query(obj, sqlstring)
	}
	this.execute(obj, id, mt_find, fun)
	return result
}

func (this *MapperDao) Insert(obj utils.Object, id string) (int64, error) {

	return 0, nil
}

func (this *MapperDao) Update(obj utils.Object, id string) (int64, error) {
	return 0, nil
}

func (this *MapperDao) Delete(obj utils.Object, id string) (int64, error) {
	return 0, nil
}

func (this *MapperDao) Query(obj utils.Object, page Page, id string) []interface{} {
	return nil
}

func (this *MapperDao) Count(id string) int64 {
	return 0
}

func (this *MapperDao) Exist(id string) bool {
	return false
}
