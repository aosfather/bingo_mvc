package hippo

import (
	"fmt"
	"strings"
)

//字段类型
type FType byte

//类型定义
const (
	URL FType = iota
	STRING
	INT
	DOUBLE
	BOOLEAN
	DATE
	DATETIME
)

func (this FType) ToString() string {
	switch this {
	case URL:
		return "U"
	case STRING:
		return "S"
	case INT:
		return "I"
	case DOUBLE:
		return "D"
	case BOOLEAN:
		return "B"
	case DATE:
		return "C"
	case DATETIME:
		return "T"
	}
	return "UNDEFINED!"
}
func ToFieldType(s string) FType {
	s = strings.ToUpper(s)
	switch s {
	case "S":
		return STRING
	case "I":
		return INT
	case "D":
		return DOUBLE
	case "B":
		return BOOLEAN
	case "C":
		return DATE
	case "T":
		return DATETIME
	case "U":
		return URL
	default:
		return STRING
	}

}

func (this *FType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var text string
	unmarshal(&text)
	*this = ToFieldType(text)
	return nil
}

//VT 值的类型
type ValueType string

const (
	VT_SINGLES   ValueType = "S" //单个值
	VT_SET       ValueType = "T" //集合或者叫枚举
	VT_RANGE     ValueType = "R" //区间 [起始值;结束值]
	VT_REGEX     ValueType = "X" //正则表达式，复合正则表达式的允许
	VT_ALL       ValueType = "*" //允许所有值
	VT_MICRO     ValueType = "&" //宏，需要解释。例如 本部门、上级部门等
	VT_STARTWITH ValueType = "B" //以指定的值开头
	VT_ENDWITH   ValueType = "E" //以指定的值为结尾
)

//执行策略
type Policy byte

const (
	P_DENY    Policy = 1 //拒绝
	P_FILTER  Policy = 2 //过滤
	P_OBSCURE Policy = 3 //混淆
)

//元信息
//权限字段
type AuthField struct {
	Code  string //编码、字段名称
	Label string //名称
	Type  FType
}

//权限表定义
type AuthTable struct {
	Code       string //编码
	Label      string //名称
	Fields     []AuthField
	AuthPolicy Policy //执行策略，1、拒绝 2、过滤 3、混淆
}

func (this *AuthTable) AddField(f AuthField) {
	if f.Code != "" {
		this.Fields = append(this.Fields, f)
	}

}
func (this *AuthTable) GetField(fname string) AuthField {
	for _, f := range this.Fields {
		if f.Code == fname {
			return f
		}
	}
	return AuthField{}
}

func (this *AuthTable) Add(code string, label string, t FType) {
	this.AddField(AuthField{code, label, t})
}

//权限列取值
type AuthColumn struct {
	Field    *AuthField //对应的字段
	Type     ValueType  //值类型
	Value    string     //取值
	Negation bool       //取反
}

func (this *AuthColumn) Validate(obj interface{}) bool {
	str := fmt.Sprintf("%v", obj)
	var typeFlag = this.Type
	var result bool = false
	switch typeFlag {
	case VT_SINGLES:
		result = singleEquals(str, this)
	case VT_RANGE:
		result = isInRange()
	case VT_SET:
		result = isInSet(str, this)
	case VT_STARTWITH:
		result = isStartWith(str, this)
	case VT_ENDWITH:
		result = isEndWith(str, this)
	case VT_REGEX:
		result = isMatchRegex(str, this)
	case VT_MICRO:
		result = isEqualsMicro()
	case VT_ALL:
		result = true
	default:
		result = false

	}
	if this.Negation {
		return !result
	}

	return result
}

//权限数据
type AuthRow map[string]*AuthColumn //所有列的值
func (this *AuthRow) Add(f *AuthField, t ValueType, v string, negation bool) {
	if f == nil || v == "" {
		return
	}
	col := &AuthColumn{Field: f, Type: t, Value: v, Negation: negation}
	row := *this
	row[f.Code] = col
}
func (this *AuthRow) AddColumn(col *AuthColumn) {
	if col != nil {
		row := *this
		row[col.Field.Code] = col
	}
}
func (this *AuthRow) HasPermition(parameters map[string]interface{}) bool {
	cols := *this
	for key, value := range cols {
		if !value.Validate(parameters[key]) {
			return false
		}
	}

	return true
}

//权限数据集合
type AuthCollection struct {
	Table    string
	rows     []*AuthRow //正向授权
	vetoRows []*AuthRow //否决
}

func (this *AuthCollection) AddAuth(r *AuthRow) {
	if r != nil {
		this.rows = append(this.rows, r)
	}

}

func (this *AuthCollection) AddVeto(r *AuthRow) {
	if r != nil {
		this.vetoRows = append(this.vetoRows, r)
	}
}

/**
判断在权限集合中是否有该权限
parameters 参数
collection 权限集合
*/

func (this *AuthCollection) HasPermition(parameters map[string]interface{}) bool {
	result := false
	//white list
	for _, value := range this.rows {
		if value.HasPermition(parameters) {
			result = true
			break
		}
	}

	//当有权限的时候，检查否决列表，如果有拒绝，则被否决
	if result {
		if this.vetoRows == nil || len(this.vetoRows) == 0 {
			return result
		}

		//black list
		for _, value := range this.vetoRows {
			if value.HasPermition(parameters) {
				return false
			}
		}

		return true
	}

	return result

}
