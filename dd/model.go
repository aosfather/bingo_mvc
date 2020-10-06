package dd

import (
	"github.com/aosfather/bingo_utils"
	"strconv"
	"strings"
)
/**
  基本模型
   1、base type
      字符串
      数字
      枚举
      日期
   2、领域元素
      长度，通用规则
   3、数据元素
      可以继承领域元素，设置check，convert

 */
type Type byte
const (
	T_string Type=iota
	T_int
	T_date
	T_datetime
)
func(this *Type) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	text=strings.ToLower(text)
	switch text {
	case "string":
		*this=T_string
	case "int":
		*this=T_int
	case "date":
		*this=T_date
	case "datetime":
		*this=T_datetime
	default:
		*this=T_string
	}
	return nil
}
type DomainType struct {
	Code string
	Label string
	DataType Type  `yaml:"type"`//数据元素
	Length   int   //最大长度
	NullEnabled bool `yaml:"null"`//是否允许为空
	Validater validatefunction //校验器
	ValidateExpr      string `yaml:"expr"`//校验表达式
}

//返回，是否校验通过及对应的错误码
func(this *DomainType)Validate(v interface{}) (bool,int){
	if this.Validater!=nil {
		b,msg:=this.Validater(this.Code,this.ValidateExpr,v)
		if !b{
			bingo_utils.Debug(msg)
			return false,103
		}
	}
	if (this.DataType==T_string) {
		s:=v.(string)
        return this.validateStr(s)
	}else if (this.DataType==T_int) {
		s:=strconv.Itoa(v.(int))
		return this.validateStr(s)
	}
	return true,0
}

//检查字符串
func (this *DomainType)validateStr(s string)(bool,int){
	if !this.NullEnabled {
		if s=="" {
			return false,102
		}
	}
	if len(s)<=this.Length {
		return true,0
	}else {
		return false,101
	}
}

//数据元素
type DataElement struct {
	Code string
	Label string
	DataType string  //所属的域元素
	domain *DomainType
	Validater validatefunction //校验器
	ValidateExpr      string //校验表达式
	Convertor  string //转换器
	ConvertorExpr string
	_convertor Convertor
}

func (this *DataElement)Input(v interface{}) interface{}{
	if v!=nil {
		if this._convertor!=nil {
			return this._convertor.Input(v)
		}
		if this.Convertor!=""{
			if factory,ok:=_convertors[this.Convertor];ok{
				this._convertor=factory(this.ConvertorExpr)
				return this._convertor.Input(v)
			}
		}
	}
	return nil
}

func (this *DataElement)Output(v interface{}) interface{}{
	if v!=nil {
		if this._convertor!=nil {
			return this._convertor.Output(v)
		}
		if this.Convertor!=""{
			if factory,ok:=_convertors[this.Convertor];ok{
				this._convertor=factory(this.ConvertorExpr)
				return this._convertor.Output(v)
			}
		}
	}
	return nil
}

func (this *DataElement)Validate(v interface{}) (bool,int){
	//使用校验器校验，通过后则进行domain校验
	if this.Validater!=nil {
		b,msg:=this.Validater(this.Code,this.ValidateExpr,v)
		if !b{
			bingo_utils.Debug(msg)
			return false,103
		}
	}
	//使用校验器校验，如果通过则使用domain进行校验
	return this.domain.Validate(v)
}

//词条
type DictCatalog struct {
	Code  string
	Label string
	Tip   string
	Items []DictCatalogItem
}

//词条的选择项
type DictCatalogItem struct {
	Code    string            //值
	Label   string            //显示值
	Tip     string            //提示
	Virtual bool              //是否虚拟,表示存在有同名的词条
	Ord     int               //显示次序
	Extends map[string]string //扩展属性
}

//转换器工厂
type ConvertorFactory func(expr string) Convertor
//转换器
type Convertor interface {
	//输入转换
	Input(v interface{})interface{}
	//输出转换
	Output(v interface{})interface{}
}
//系统转换器
var _convertors map[string]ConvertorFactory=make(map[string]ConvertorFactory)
func RegisterConvertorFactory(name string ,c ConvertorFactory){
	if name!=""&& c!=nil {
		if _,ok:=_convertors[name];!ok{
			_convertors[name]=c
		}
	}
}