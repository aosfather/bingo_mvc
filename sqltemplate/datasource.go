package sqltemplate

import (
	"database/sql"
	"fmt"
	"strings"

	"log"
)

func debug(msg string, obj ...interface{}) {
	log.Printf(msg, obj...)
}

/**
  data source

*/
type DataSource struct {
	DBtype      string
	DBurl       string
	DBuser      string
	DBpassword  string
	DBname      string
	pool        *sql.DB
	sqlTemplate *SqlTemplate
}

func (this *DataSource) Init() {
	this.sqlTemplate = &SqlTemplate{}

	//如果已经初始化，不在初始化
	if this.pool != nil {
		return
	}

	if strings.ToLower(this.DBtype) == "mysql" {
		dburl := this.DBurl
		if strings.Index(dburl, "(") <= 0 {
			dburl = fmt.Sprintf("tcp(%s)", dburl)
		}
		url := this.DBuser + ":" + this.DBpassword + "@" + dburl + "/" + this.DBname
		var err error
		this.pool, err = sql.Open(this.DBtype, url)
		if err == nil {
			this.pool.Ping()
		} else {
			debug("%v", err)
		}

	}
}

//获取连接
func (this *DataSource) GetConnection() *Connection {
	var conn Connection
	conn.db = this.pool
	conn.template = this.sqlTemplate
	return &conn
}

func (this *DataSource) GetDao() *BaseDao {
	return &BaseDao{this}
}
