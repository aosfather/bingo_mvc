package bingo_dao

import (
	"fmt"
	utils "github.com/aosfather/bingo_utils"
)

/**
  DAO 面向结构进行自动sql生成
*/
const _maxsize = 10000 //默认分页中最大条数

//基础数据操作对象
type BaseDao struct {
	ds *DataSource
}

func (this *BaseDao) Init(c *DataSource) {
	this.ds = c
}

//插入，返回auto id和错误信息
func (this *BaseDao) Insert(obj utils.Object) (int64, error) {
	if obj == nil {
		return 0, fmt.Errorf("nil object!")
	}

	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()
		id, count, err := session.Insert(obj)

		if err == nil && count == 1 {
			session.Commit()
			return id, nil
		} else {
			session.Rollback()
			return 0, err
		}

	}

	return 0, fmt.Errorf("session is nil")

}

//插入，返回auto id和错误信息
func (this *BaseDao) InsertBatch(objs []utils.Object) (ids []int64, e error) {
	if objs == nil || len(objs) == 0 {
		return nil, fmt.Errorf("nil object!")
	}

	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()

		for _, obj := range objs {
			id, _, err := session.Insert(obj)

			if err != nil {
				session.Rollback()
				return nil, err
			}
			ids = append(ids, id)
		}

		session.Commit()
		e = nil
		return

	}

	return nil, fmt.Errorf("session is nil")

}

func (this *BaseDao) Find(obj utils.Object, cols ...string) bool {
	if obj == nil {
		return false
	}
	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		return session.Find(obj, cols...)
	}
	return false
}

func (this *BaseDao) FindBySql(obj utils.Object, sqlTemplate string, args ...interface{}) bool {
	if obj == nil {
		return false
	}
	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		return session.Query(obj, sqlTemplate, args...)
	}
	return false
}

//更新，返回更新的条数和错误信息
func (this *BaseDao) Update(obj utils.Object, cols ...string) (int64, error) {
	if obj == nil {
		return 0, fmt.Errorf("nil object!")
	}

	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()
		_, count, err := session.Update(obj, cols...)
		if err != nil {
			session.Rollback()
			return 0, err
		} else {
			session.Commit()
			return count, nil
		}
	}

	return 0, fmt.Errorf("session is nil")
}

func (this *BaseDao) UpdateBatch(objs []utils.Object, cols ...string) (counts []int64, e error) {
	if objs == nil || len(objs) == 0 {
		return nil, fmt.Errorf("nil object!")
	}

	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()
		for _, obj := range objs {
			_, count, err := session.Update(obj, cols...)
			if err != nil {
				session.Rollback()
				return nil, err
			}
			counts = append(counts, count)

		}
		session.Commit()
		e = nil
		return
	}

	return nil, fmt.Errorf("session is nil")
}

func (this *BaseDao) Delete(obj utils.Object, cols ...string) (int64, error) {
	if obj == nil {
		return 0, fmt.Errorf("nil object!")
	}

	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()
		_, count, err := session.Delete(obj, cols...)
		if err != nil {
			session.Rollback()
			return 0, err
		} else {
			session.Commit()
			return count, nil
		}
	}

	return 0, fmt.Errorf("session is nil")
}

func (this *BaseDao) QueryAll(obj utils.Object, cols ...string) []interface{} {

	page := Page{_maxsize, 0, 0}
	return this.Query(obj, page, cols...)

}

func (this *BaseDao) QueryAllBySql(obj utils.Object, sqlTemplate string, args ...interface{}) []interface{} {

	page := Page{_maxsize, 0, 0}
	if obj == nil {
		return nil
	}
	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		return session.QueryByPage(obj, page, sqlTemplate, args...)

	}

	return nil

}

func (this *BaseDao) Query(obj utils.Object, page Page, cols ...string) []interface{} {

	if obj == nil {
		return nil
	}
	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		theSql, args, err := this.ds.sqlTemplate.CreateQuerySql(obj, cols...)
		if err == nil {
			return session.QueryByPage(obj, page, theSql, args...)
		}

	}

	return nil
}

//执行单条sql
func (this *BaseDao) Exec(sqltemplate string, objs ...interface{}) (int64, error) {
	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()
		_, count, err := session.ExeSql(sqltemplate, objs...)
		if err == nil {
			session.Commit()
			return count, err
		}

		session.Rollback()
		return 0, err
	}
	return 0, fmt.Errorf("session is nil")
}

//批量执行简单的sql语句
func (this *BaseDao) ExecSqlBatch(sqls ...string) error {
	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()
		for _, sql := range sqls {
			_, _, err := session.ExeSql(sql)
			if err != nil {
				session.Rollback()
				return err
			}

		}
		session.Commit()
		return nil
	}
	return fmt.Errorf("session is nil")
}

func (this *BaseDao) GetSession() *Connection {
	return this.ds.GetConnection()
}

//插入对象，并根据对象的id，更新后续的sql语句，一般为update，其中这个关联id必须是第一个参数
func (this *BaseDao) InsertAndUpdate(iobj interface{}, sqltemplate string, args ...interface{}) error {
	if iobj == nil {
		return fmt.Errorf("nil object!")
	}

	session := this.ds.GetConnection()
	if session != nil {
		defer session.Close()
		session.Begin()
		id, count, err := session.Insert(iobj)
		if err == nil && count == 1 {
			p := []interface{}{id}
			if args != nil && len(args) > 0 {
				p = append(p, args...)
			}

			_, _, err = session.ExeSql(sqltemplate, p)

			if err == nil {
				session.Commit()
			}

		}

		session.Rollback()
		return err
	}

	return fmt.Errorf("session is nil")

}
