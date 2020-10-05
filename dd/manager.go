package dd

import (
	"fmt"
	files "github.com/aosfather/bingo_utils/files"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)
var _domains map[string]*DomainType=make(map[string]*DomainType)
var _elements map[string]*DataElement=make(map[string]*DataElement)
func AddDomainType (domainType DomainType) error {
	code:=domainType.Code
	if isCodeExist(code) {
		return fmt.Errorf("the domain type name [%s] is exist!",code)
	}
	_domains[code]=&domainType
	return nil
}

func AddDataElement(d DataElement)error {
	code:=d.Code
	if isCodeExist(code) {
		return fmt.Errorf("the domain type name [%s] is exist!",code)
	}
	d.domain=getDomainType(d.DataType)
	_elements[code]=&d
	return nil
}

func isCodeExist(code string)bool {
	if _,ok:=_domains[code];ok {
		return true
	}

	if _,ok:=_elements[code];ok {
		return true
	}

	return false
}
func getDomainType(n string) *DomainType {
	return _domains[n]
}
func getDataElement(n string)*DataElement{
	return _elements[n]
}

type configFile struct {
    Domains []DomainType
    Elements []DataElement
    Dicts []DictCatalog
}

//加载domain和dataelement 定义
func LoadConfig(cf string){
	if files.IsFileExist(cf){
		config:=configFile{}
		data, err := ioutil.ReadFile(cf)
		if err == nil {
			err = yaml.Unmarshal(data, &config)
		}
		if err != nil {
			//errs("load verify meta error:", err.Error())
			return
		}
        //加载领域定义
		for _, item := range config.Domains {
			AddDomainType(item)
		}
        //加载数据元素定义
		for _, item := range config.Elements {
			AddDataElement(item)
		}
		//加载数据字典
		for _, item := range config.Dicts {
			_dicts[item.Code]=&item
		}
		SetDictMeta(getDictCatalog)

	}
}

//默认dict meta实现
var _dicts map[string]*DictCatalog=make(map[string]*DictCatalog)
func getDictCatalog(d string)*DictCatalog{
	return _dicts[d]
}