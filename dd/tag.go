package dd

import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
	utils "github.com/aosfather/bingo_utils"
	reflect2 "github.com/aosfather/bingo_utils/reflect"
	"reflect"
)
func init(){
	bingo_mvc.RegisterValidate(Validate)
}
const (
	_TAG_DD = "DD" //属性
)
var DefaultErorrFormat="The %s field validate failed by '%s'"
func Validate(obj interface{}) []error {
	if reflect2.IsMap(obj) {
		//TODO 需要看怎么做校验了
		return nil
	}
	objT, objV, err := reflect2.GetStructTypeValue(obj)
	if err != nil {
		return []error{utils.CreateError(501, err.Error())}
	}
	var errors []error
	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		tag := f.Tag.Get(_TAG_DD)
		v := vf.Interface()
		if len(tag) == 0 {
			if vf.Kind() == reflect.Struct || (vf.Kind() == reflect.Ptr && vf.Elem().Kind() == reflect.Struct) {
				errors = append(errors, Validate(v)...)
			}

		} else {
			e:=ValidateByDataType(f.Name,v,tag)
			if e!=nil {
				errors = append(errors, e)
			}

		}

	}
	return errors

}

//使用指定的数据元素对值进行校验
func ValidateByDataType(name string,v interface{},dt string) error{
	de:=getDataElement(dt)
	if de!=nil {
		b,e:= de.Validate(v)
		if !b{
			return fmt.Errorf(DefaultErorrFormat,name,GetErrorMsg(e))
		}
	}
	//查找对应的域对象
	domain:=getDomainType(dt)
	if domain!=nil {
		b,e:=domain.Validate(v)
		if !b{
			return fmt.Errorf(DefaultErorrFormat,name,GetErrorMsg(e))
		}
	}

	return nil
}
