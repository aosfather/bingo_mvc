package sqltemplate

import (
	"database/sql"
	"fmt"
	"github.com/aosfather/bingo_utils/files"
	"io/ioutil"
	"os"
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
	DBtype             string
	DBurl              string
	DBuser             string
	DBpassword         string
	DBname             string
	DBmapper           string //mapper文件夹
	pool               *sql.DB
	sqlTemplate        *SqlTemplate
	sqlTemplateManager *SqltemplateManager
}

func (this *DataSource) Init() {
	this.sqlTemplate = &SqlTemplate{}
	//构建sqltemplatemanager
	if this.DBmapper != "" {
		if files.IsFileExist(this.DBmapper) {
			this.sqlTemplateManager = &SqltemplateManager{}
			this.loadmapper(string(os.PathSeparator), this.DBmapper)
		}
	}

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

func (this *DataSource) loadmapper(pathSeparator string, fileDir string) {
	files, _ := ioutil.ReadDir(fileDir)
	for _, onefile := range files {
		filename := fileDir + pathSeparator + onefile.Name()
		if onefile.IsDir() {
			this.loadmapper(pathSeparator, filename)
		} else {
			this.sqlTemplateManager.AddCollectFromFile(filename)
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

func (this *DataSource) GetMapperDao(namespace string) *MapperDao {
	if this.sqlTemplateManager != nil {
		return this.sqlTemplateManager.BuildDao(this, namespace)
	}
	return nil
}
