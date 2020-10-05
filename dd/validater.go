package dd

import (
	"fmt"
	"github.com/aosfather/bingo_utils"
	"regexp"
)

//校验函数
type validatefunction func(code string,expr string ,v interface{}) (bool,string)
func(this *validatefunction) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	*this = _validatefunctions[text]
	return nil
}

var _validatefunctions map[string]validatefunction=make(map[string]validatefunction)
//注册函数
func registerValidateFunction(name string ,vf validatefunction) {
	if name!=""||vf==nil {
		return
	}
	if _,ok:=_validatefunctions[name];ok{
		bingo_utils.Debug(name," ,validate function exist!")
		return
	}
	_validatefunctions[name]=vf
}

//----------------------正则表达式-----------------------------//
var regexpMap map[string]*regexp.Regexp=make(map[string]*regexp.Regexp)
func validateByRegexp(code string,expr string ,v interface{}) (bool, string) {
	var pattern *regexp.Regexp
	var err error
    if pattern,ok:=regexpMap[code];!ok{
		pattern, err = regexp.Compile(expr)
		if err!=nil{
			bingo_utils.Err(err.Error())
		}else {
			regexpMap[code]=pattern
		}
	}

	if pattern != nil {
		result := pattern.Match([]byte(v.(string)))
		if result {
			return true, ""
		}
	}

	return false, "regex校验不同过！"

}
//---------------------------字典------------------------//
type DictMeta func(code string) *DictCatalog
var _dictMeta DictMeta
func SetDictMeta(dm DictMeta){
	if dm!=nil {
		_dictMeta=dm
	}

}

func GetDict(d string)*DictCatalog{
	return _dictMeta(d)
}
//校验字典的值
func validateByDict(code string,expr string ,v interface{}) (bool, string) {
	if _dictMeta==nil {
		return true,"no dict meta manager!"
	}
	if c := _dictMeta(expr); c != nil {
		for _, item := range c.Items {
			if item.Code == v || item.Label == v {
				return true, ""
			}
		}
		return false, fmt.Sprintf("'%s'不是字典[%s]的合法成员", v, expr)

	}
	return false, fmt.Sprintf("指定的字典[%s],不存在", expr)
}
//----------------------------------------------------------//
func init(){
	registerValidateFunction("regexp",validateByRegexp)
	registerValidateFunction("dict",validateByDict)
}

