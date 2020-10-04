package hippo

/*
角色
1、由权限集合构成
2、可以继承角色，扩充被继承的角色
3、可以设置管控角色，角色的权限受制于管控角色
角色     管控角色     结果
 T		  T			 T
 F		  T			 F
 T        F         F
 F        F         F
*/
type Role struct {
	Code          string //编码、字段名称
	Label         string //名称
	ParentRole    *Role  //父类角色
	ManagedRole   *Role  //受管理角色
	collectionMap map[string]*AuthCollection
}

func (this *Role) Init() {
	this.collectionMap = make(map[string]*AuthCollection)
}

//按权限表加入权限数据
func (this *Role) AddRow(table string, row *AuthRow, veto bool) {
	if row != nil && table != "" {
		var c *AuthCollection
		c = this.collectionMap[table]
		if c == nil {
			c = &AuthCollection{Table: table}
			this.collectionMap[table] = c
		}
		if veto {
			c.AddVeto(row)
		} else {
			c.AddAuth(row)
		}
	}
}

func (this *Role) HasPermition(tableTrigger string, parameters map[string]interface{}) bool {
	if ac, ok := this.collectionMap[tableTrigger]; ok {
		return ac.HasPermition(parameters)
	}
	return false
}
