/**
SQL生成模板
用于根据struct的定义生成update语句和insert语句
规则
  1、表名
   默认t_stuct的名称(小写，驼峰转 xx_xx)
   可以通过setTablePreFix来指定默认的表前缀
   可以通过tag Table:""指定表名

  2、字段名称
    默认字段名称(小写，驼峰转 xx_xx)
	可以通过 tag Field:""指定字段名

  3、特例
    通过tag Option:"" 来指定。
	可选值：auto、pk、not 分别表示 自动增长、主健、忽略

*/
package sqltemplate

import (
	"fmt"
	reflect2 "github.com/aosfather/bingo_utils/reflect"
	strings2 "github.com/aosfather/bingo_utils/strings"
	"reflect"
	"strings"
)

const (
	_TAG_FIELD   = "Field"
	_Tag_table   = "Table"
	table_prefix = "t_"
)

var (
	table_insert_cache map[string]string
)

func init() {
	table_insert_cache = make(map[string]string)
}

type SqlTemplate struct {
	table_prefix string
}

func (this *SqlTemplate) SetTablePrefix(prefix string) {
	this.table_prefix = prefix
}

func (this *SqlTemplate) GetInsertSql(target interface{}) (string, []interface{}, error) {
	objT, _, err := reflect2.GetStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	key := objT.Name()
	sql := table_insert_cache[key]

	if sql == "" {
		sql, args, err := this.CreateInserSql(target)
		if err != nil {
			return "", nil, err
		}
		table_insert_cache[key] = sql
		return sql, args, err

	}

	args, err := this.structValueToArray(target)
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func (this *SqlTemplate) StructValueToCustomArray(target interface{}, col ...string) ([]interface{}, error) {
	_, objV, err := reflect2.GetStructTypeValue(target)
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, len(col))
	for i, field := range col {
		vf := objV.FieldByName(field)
		if !vf.CanInterface() {
			args[i] = nil
		} else {
			args[i] = objV.FieldByName(field).Interface()
		}

	}
	return args, nil
}

func (this *SqlTemplate) structValueToArray(target interface{}) ([]interface{}, error) {
	objT, objV, err := reflect2.GetStructTypeValue(target)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0, 0)
	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		//对于自增长和明确忽略的字段不做转换
		if this.isFieldIgnore(f) {
			continue
		}

		args = append(args, vf.Interface())

	}
	return args, nil

}

func (this *SqlTemplate) CreateQuerySql(target interface{}, col ...string) (string, []interface{}, error) {
	objT, objV, err := reflect2.GetStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	var tagTableName string
	var sqlFields, sqlwheres []string
	var argsWhere, args []interface{}

	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		//处理内嵌结构
		if f.Anonymous {
			this.addEmberStruct(f, vf, "", &sqlFields, &sqlwheres, &args, &argsWhere)
			continue
		}

		tagTable := f.Tag.Get(_Tag_table)
		if tagTable != "" {
			tagTableName = tagTable
		}

		this.addFieldAndWhere(f, vf, "", &sqlFields, &sqlwheres, &args, &argsWhere)

	}

	//如果没有指定表名就使用默认规则
	if tagTableName == "" {
		tagTableName = this.getDefaultTableName(objT)
	}

	return "select " + strings.Join(sqlFields, ",") + " from " + tagTableName + " where " + strings.Join(sqlwheres, ","), argsWhere, nil
}

func (this *SqlTemplate) CreateDeleteSql(target interface{}, col ...string) (string, []interface{}, error) {
	objT, objV, err := reflect2.GetStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	var tagTableName string
	var sqlFields, sqlwheres []string
	var argsWhere, args []interface{}

	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		//处理内嵌结构
		if f.Anonymous {
			this.addEmberStruct(f, vf, "", &sqlFields, &sqlwheres, &args, &argsWhere)
			continue
		}

		tagTable := f.Tag.Get("Table")
		if tagTable != "" {
			tagTableName = tagTable
		}

		this.addFieldAndWhere(f, vf, "", &sqlFields, &sqlwheres, &args, &argsWhere)

	}

	//如果没有指定表名就使用默认规则
	if tagTableName == "" {
		tagTableName = this.getDefaultTableName(objT)
	}

	return "delete from " + tagTableName + " where " + strings.Join(sqlwheres, " and "), argsWhere, nil
}

func (this *SqlTemplate) CreateInserSql(target interface{}) (string, []interface{}, error) {
	objT, objV, err := reflect2.GetStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	var tagTableName string
	var sqlFields []string
	var sqlValues []string
	var args []interface{}

	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}
		tagTable := f.Tag.Get("Table")
		if tagTable != "" {
			tagTableName = tagTable
		}

		//对于自增长和明确忽略的字段不做转换
		tagOption := f.Tag.Get("Option")
		if tagOption != "" {
			if strings.Index(tagOption, "auto") != -1 || strings.Index(tagOption, "not") != -1 {
				continue
			}
		}

		this.addFields(f, vf, &sqlFields, &sqlValues, &args)

	}

	//如果没有指定表名就使用默认规则
	if tagTableName == "" {
		tagTableName = this.getDefaultTableName(objT)
	}

	return "Insert into " + tagTableName + "(" + strings.Join(sqlFields, ",") + ") Values(" + strings.Join(sqlValues, ",") + ")", args, nil
}

func (this *SqlTemplate) CreateUpdateSql(target interface{}, col ...string) (string, []interface{}, error) {
	objT, objV, err := reflect2.GetStructTypeValue(target)
	if err != nil {
		return "", nil, err
	}
	var tagTableName string
	var sqlwheres, sqlFields []string

	args := make([]interface{}, 0, 0)
	argsWhere := make([]interface{}, 0, 0)

	//遍历字段
	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		//处理内嵌结构
		if f.Anonymous {
			this.addEmberStruct(f, vf, "=?", &sqlFields, &sqlwheres, &args, &argsWhere)
			continue
		}

		//获取指明的table名称
		tagTable := f.Tag.Get("Table")
		if tagTable != "" {
			tagTableName = tagTable
		}

		this.addFieldAndWhere(f, vf, "=?", &sqlFields, &sqlwheres, &args, &argsWhere)

	}

	//如果没有指定表名就使用默认规则
	if tagTableName == "" {
		tagTableName = this.getDefaultTableName(objT)
	}

	args = append(args, argsWhere...)
	return "update " + tagTableName + " set " + strings.Join(sqlFields, ",") + " where " + strings.Join(sqlwheres, " and "), args, nil

}

func (this *SqlTemplate) addEmberStruct(f reflect.StructField, v reflect.Value, fieldfix string, fields, wheres *[]string, args, argsWhere *[]interface{}) {
	if f.Anonymous {
		ft := f.Type
		for i := 0; i < ft.NumField(); i++ {
			ff := ft.Field(i)
			vf := v.Field(i)
			if !vf.CanInterface() {
				continue
			}
			//处理内嵌结构
			if ff.Anonymous {
				this.addEmberStruct(ff, vf, fieldfix, fields, wheres, args, argsWhere)
				continue
			}

			this.addFieldAndWhere(ff, vf, fieldfix, fields, wheres, args, argsWhere)
		}
	}
}

func (this *SqlTemplate) addFields(f reflect.StructField, v reflect.Value, fields, values *[]string, args *[]interface{}) {
	if f.Anonymous {
		fmt.Println("em")
		ft := f.Type
		for i := 0; i < ft.NumField(); i++ {
			ff := ft.Field(i)
			vf := v.Field(i)
			if !vf.CanInterface() {
				continue
			}
			//处理内嵌结构
			if ff.Anonymous {
				this.addFields(ff, vf, fields, values, args)
				continue
			}

			this.addField(ff, vf, fields, values, args)
		}
	} else {
		this.addField(f, v, fields, values, args)
	}
}

func (this *SqlTemplate) addField(f reflect.StructField, v reflect.Value, fields, values *[]string, args *[]interface{}) {
	colName := GetColName(f)
	fmt.Println(colName)
	fmt.Println(v.Interface())
	*args = append(*args, v.Interface())
	*fields = append(*fields, colName)
	*values = append(*values, "?")
}

func (this *SqlTemplate) addFieldAndWhere(f reflect.StructField, v reflect.Value, fieldfix string, fields, wheres *[]string, args, argsWhere *[]interface{}) {
	colName := reflect2.GetColName(f)
	//
	*fields = append(*fields, colName+fieldfix)
	*args = append(*args, v.Interface())

	//对于自增长和明确忽略的字段不做转换
	tagOption := f.Tag.Get("Option")
	if tagOption != "" {
		if strings.Index(tagOption, "pk") != -1 {
			//where的处理
			*wheres = append(*wheres, colName+" =?")
			*argsWhere = append(*argsWhere, v.Interface())
		}
	}

}

//获取默认的表名
func (this *SqlTemplate) getDefaultTableName(t reflect.Type) string {
	if this.table_prefix == "" {
		this.table_prefix = table_prefix
	}

	return this.table_prefix + strings2.BingoString(t.Name()).SnakeString()
}

//字段是否忽略
func (this *SqlTemplate) isFieldIgnore(field reflect.StructField) bool {
	//对于自增长和明确忽略的字段不做转换
	tagOption := field.Tag.Get("Option")
	if tagOption != "" {
		if strings.Index(tagOption, "auto") != -1 || strings.Index(tagOption, "not") != -1 {
			return true
		}
	}
	return false
}

func GetColName(field reflect.StructField) string {
	colName := field.Tag.Get(_TAG_FIELD)
	if colName == "" {
		colName = strings2.BingoString(field.Name).SnakeString()
	}

	return colName
}
