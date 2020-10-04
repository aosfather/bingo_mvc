package hippo

import (
	"github.com/aosfather/bingo_utils/files"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//原信息接口
type TableMetaReader interface {
	GetTable(string) *AuthTable
}

//角色信息读取器
type RoleReader interface {
	GetRole(string) *Role
}

//权限引擎
type AuthEngine struct {
	TableMeta TableMetaReader `Inject:""`
	RoleMeta  RoleReader      `Inject:""`
}

/*
tableTrigger 触发标识，用于识别出对应的权限表
parameters   参数
roles 角色列表
*/
func (this *AuthEngine) HasPermition(tableTrigger string, parameters map[string]interface{}, roles ...string) bool {
	for _, role := range roles {
		roleObj := this.RoleMeta.GetRole(role)
		if roleObj != nil {
			if roleObj.HasPermition(tableTrigger, parameters) {
				return true
			}
		}

	}
	return false
}

type TableFile struct {
	Version     string
	Description string
	Tables      []AuthTable
}
type YamlFileTableMeta struct {
	tables map[string]*AuthTable
}

func (this *YamlFileTableMeta) Load(f string) {
	this.tables = make(map[string]*AuthTable)
	if files.IsFileExist(f) {
		tf := &TableFile{}
		data, err := ioutil.ReadFile(f)
		if err == nil {
			err = yaml.Unmarshal(data, tf)
		}
		if err != nil {
			//errs("load verify meta error:", err.Error())
			return
		}

		for _, item := range tf.Tables {
			this.tables[item.Code] = &item

		}
	}
}
func (this *YamlFileTableMeta) GetTable(table string) *AuthTable {
	if this.tables != nil {
		return this.tables[table]
	}

	return nil
}
