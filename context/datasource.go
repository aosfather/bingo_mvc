package context

import (
	"github.com/aosfather/bingo_mvc/sqltemplate"
	"github.com/aosfather/bingo_utils/reflect"
)

func InitDatasource(f reflect.StoreFunction) (string, interface{}) {
	if f("bingo.usedb") == "true" {
		var ds sqltemplate.DataSource
		ds.DBtype = f("bingo.db.type")
		ds.DBname = f("bingo.db.name")
		ds.DBurl = f("bingo.db.url")
		ds.DBuser = f("bingo.db.user")
		ds.DBpassword = f("bingo.db.password")
		return "", &ds
	}
	return "", nil

}
