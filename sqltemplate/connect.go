package sqltemplate

import (
	"database/sql"
	"fmt"
	utils "github.com/aosfather/bingo_utils"
	"github.com/aosfather/bingo_utils/reflect"
)

//最大重试次数
const Max_RETRY = 10

/**
  数据库 连接包装
*/
type Connection struct {
	tx       *sql.Tx
	db       *sql.DB
	isTx     bool //是否开启了事务
	template *SqlTemplate
}

func (this *Connection) Begin() {
	if this.isTx {
		utils.Err("tx has opened!")
		return
	}
	times := 0
opentx:
	var err error
	this.tx, err = this.db.Begin()
	if err != nil {
		if times > Max_RETRY {
			utils.Err("retry times reach max times! ")
			return
		}

		times++
		utils.Err("db open error ", err.Error(), " retry the ", times, " times!")

		goto opentx
	}
	this.isTx = true

}

func (this *Connection) Commit() {
	if this.tx != nil && this.isTx {
		this.tx.Commit()
		this.isTx = false
	}
}

func (this *Connection) Rollback() {
	if this.tx != nil && this.isTx {
		this.tx.Rollback()
		this.isTx = false
	}
}

func (this *Connection) prepare(sql string) (*sql.Stmt, error) {
	if this.isTx {
		//utils.Debugf("%v", this.tx)
		return this.tx.Prepare(sql)
	} else if this.db != nil {
		return this.db.Prepare(sql)
	}
	return nil, utils.CreateError(500, "no db init")
}

func (this *Connection) Close() {
	this.Rollback()
}

func (this *Connection) SimpleQuery(sql string, obj ...interface{}) bool {
	stmt, err := this.prepare(sql)
	if err != nil {
		utils.Debugf("%v", err)
		return false
	}
	defer stmt.Close()
	rs, err := stmt.Query()
	if err != nil {
		utils.Debugf("%v", err)
		return false
	}
	defer rs.Close()
	if rs.Next() {
		rs.Scan(obj...)
		return true
	}

	return false

}

func (this *Connection) ExeSql(sql string, objs ...interface{}) (id int64, affect int64, err error) {
	stmt, err := this.prepare(sql)
	if err != nil {
		utils.Debugf("%v", err)
		return 0, 0, err
	}
	defer stmt.Close()

	rs, err := stmt.Exec(objs...)
	if err != nil {
		utils.Debugf("%v", err)
		return 0, 0, err
	}
	id, _ = rs.LastInsertId()
	affect, _ = rs.RowsAffected()
	return id, affect, nil

}

func (this *Connection) QueryByPage(result interface{}, page Page, sql string, objs ...interface{}) []interface{} {
	//使用真分页的方式实现
	stmt, err := this.prepare(sql + buildMySqlLimitSql(page))
	if err != nil {
		utils.Debugf("%v", err)
		return nil
	}
	defer stmt.Close()
	rs, err := stmt.Query(objs...)
	if err != nil {
		utils.Debugf("%v", err)
		return nil
	}
	defer rs.Close()

	resultArray := []interface{}{}
	resultType := reflect.GetRealType(result)
	cols, _ := rs.Columns()
	for {
		if rs.Next() {
			columnsMap := make(map[string]interface{}, len(cols))
			refs := make([]interface{}, 0, len(cols))
			for _, col := range cols {
				var ref interface{}
				columnsMap[col] = &ref
				refs = append(refs, &ref)
			}

			rs.Scan(refs...)
			var arrayItem interface{}

			//填充result
			if reflect.IsMapPtr(result) || reflect.IsMap(result) {
				targetResult := make(map[string]interface{})
				this.fillToMap(columnsMap, &targetResult)
				arrayItem = targetResult
			} else {
				arrayItem = reflect.CreateObjByType(resultType)
				reflect.FillStruct(columnsMap, arrayItem)
			}

			resultArray = append(resultArray, arrayItem)

			//index++

		} else {
			break
		}
	}
	return resultArray
}
func (this *Connection) Query(result interface{}, sql string, objs ...interface{}) bool {
	stmt, err := this.prepare(sql)
	if err != nil {
		utils.Debugf("%v", err)
		return false
	}
	defer stmt.Close()
	rs, err := stmt.Query(objs...)
	if err != nil {
		utils.Debugf("%v", err)
		return false
	}
	defer rs.Close()

	//处理结构体
	if reflect.IsStructPtr(result) || reflect.IsMapPtr(result) {
		cols, _ := rs.Columns()
		columnsMap := make(map[string]interface{}, len(cols))
		refs := make([]interface{}, 0, len(cols))
		for _, col := range cols {
			var ref interface{}
			columnsMap[col] = &ref
			refs = append(refs, &ref)
		}
		if rs.Next() {
			rs.Scan(refs...)
			//填充result
			if reflect.IsMapPtr(result) {
				//填充map
				this.fillToMap(columnsMap, result)
			} else {
				reflect.FillStruct(columnsMap, result)
			}

			return true
		}
	} else { //普通指针的赋值
		if rs.Next() {
			rs.Scan(result)
			return true
		}
	}

	return false
}

//填充到map中
func (this *Connection) fillToMap(columnsMap map[string]interface{}, result interface{}) {
	target := *result.(*map[string]interface{})
	for key, value := range columnsMap {
		target[key] = reflect.GetRealValue(value)
	}
}

func (this *Connection) Insert(obj interface{}) (id int64, affect int64, err error) {
	sql, args, err := this.template.GetInsertSql(obj)
	if err != nil {
		return 0, 0, err
	}
	utils.Debugf("%v", err)
	return this.ExeSql(sql, args...)
}

func (this *Connection) Find(obj interface{}, col ...string) bool {
	sql, args, err := this.template.CreateQuerySql(obj, col...)
	if err != nil {
		return false
	}
	utils.Debugf("%v", err)
	return this.Query(obj, sql, args...)

}

func (this *Connection) Update(obj interface{}, col ...string) (id int64, affect int64, err error) {
	sql, args, err := this.template.CreateUpdateSql(obj, col...)
	if err != nil {
		return 0, 0, err
	}
	utils.Debugf("%v", sql)
	return this.ExeSql(sql, args...)
}

func (this *Connection) Delete(obj interface{}, col ...string) (id int64, affect int64, err error) {
	sql, args, err := this.template.CreateDeleteSql(obj, col...)
	if err != nil {
		return 0, 0, err
	}
	utils.Debugf("%v", sql)
	return this.ExeSql(sql, args...)
}

//mysql的分页sql生成
func buildMySqlLimitSql(page Page) string {
	if page.Index < 1 {
		page.Index = 1
	}

	if page.Size == 0 {
		page.Size = 10
	}

	return fmt.Sprintf(" limit %d,%d", page.Size*(page.Index-1), page.Size)

}
