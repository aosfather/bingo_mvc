package bingo_dao

import (
	"database/sql"
	log2 "github.com/aosfather/bingo_utils/log"
	"log"
)

var logger log2.Log

func SetLogger(l log2.Log) {
	logger = l
}

func debug(msg string, obj ...interface{}) {
	if logger != nil {
		logger.Debug(msg, obj...)
	} else {
		log.Printf(msg, obj...)
	}
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

	if this.DBtype != "" {
		url := this.DBuser + ":" + this.DBpassword + "@" + this.DBurl + "/" + this.DBname
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
