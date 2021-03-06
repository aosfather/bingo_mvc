package bingo_mvc

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aosfather/bingo_utils/files"
	reflect2 "github.com/aosfather/bingo_utils/reflect"
	"io"
	"mime"
	"os"
	"reflect"
	"strings"
)

type AutoValidateFunction func(v interface{}) []error
var _validateFunction AutoValidateFunction
func RegisterValidate(af AutoValidateFunction){
	_validateFunction=af
}
//url元信息描述
type RequestMapper struct {
	Name    string   //名称
	Url     []string //url路径
	Methods []HttpMethodType
	//请求参数类型
	Request interface{}
	//返回值处理器
	Response      Convertor
	Handle        HandleFunction
	ResponseStyle StyleType
}

func (this *RequestMapper) IsSupportMethod(m HttpMethodType) bool {
	for _, sm := range this.Methods {
		if sm == m {
			return true
		}
	}
	return false
}

func (this *RequestMapper) Select(writer io.Writer, input func(interface{}) error) StyleType {
	var paramter interface{}
	var err error
	if this.Request != nil {
		t := reflect.TypeOf(this.Request)
		if t.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
			t = t.Elem()
		}

		//获取请求参数
		if t.Kind() == reflect.Map {
			paramter = reflect.MakeMap(t).Interface()

		} else {
			paramter = reflect.New(t).Interface()
		}

		err = input(paramter)
	}

	if err == nil {
		if this.Handle != nil {
			//配置有参数自动校验的，进行自动校验
			if _validateFunction!=nil {
				errs:=_validateFunction(paramter)
				if errs!=nil&&len(errs)>0 {
					msg:=""
					for _,e:=range errs {
						msg+=e.Error()+";\n"
					}
					err=fmt.Errorf("%s",msg)
				}
			}
			//处理，并渲染
			var result interface{}
			if err!=nil{
				result = this.Handle(err)
			}else {
				result = this.Handle(paramter)
			}

			err = this.writeResponse(writer, result)
		} else {
			this.writeResponse(writer, "not found url handle!")
		}

		//错误处理进行输出
		if err != nil {
			this.writeResponse(writer, fmt.Sprintf("error:%s", err.Error()))
		}
	}

	return this.ResponseStyle

}
func (this *RequestMapper) writeResponse(w io.Writer, r interface{}) error {
	if this.Response != nil {
		return this.Response(w, r)
	} else {
		w.Write([]byte(fmt.Sprintf("%v", r)))
	}
	return nil

}

//结果集转换器
type Convertor func(writer io.Writer, obj interface{}) error

//Request响应函数
type HandleFunction func(interface{}) interface{}

//控制器
type Controller interface {
	Select(writer io.Writer, input func(interface{}) error) StyleType
	IsSupportMethod(m HttpMethodType) bool
}

//标准的map 类型
var _mapType map[string]interface{} = make(map[string]interface{})

func TypeOfMap() interface{} {
	return _mapType
}

//空请求参数，无参数请求
type _EmptyRequest struct {
}

var _Empty = new(_EmptyRequest)

func TypeOfEmpty() interface{} {
	return _Empty
}

//handle的map，用户获取handle列表
func NewHandleMap() HandleMap {
	return make(HandleMap)
}

type HandleMap map[string]CMap

func (this HandleMap) Add(name string, handle HandleFunction, p interface{}) {
	this[name] = CMap{Handle: handle, Parameter: p}
}

//controller的mapp
type CMap struct {
	Handle    HandleFunction
	Parameter interface{}
}

//多个handle的控制器
type MutiController interface {
	GetHandles() HandleMap
}

//静态资源处理
type staticControl struct {
	root string
}

func (this *staticControl) Getstaticfile(url string, writer io.Writer) (string, error) {
	index := strings.Index(url, "?")
	if index > 0 {
		url = url[0:index]
	}
	filename := this.root + url
	if files.IsFileExist(filename) {
		fixIndex := strings.LastIndex(filename, ".")
		if fixIndex > 0 {
			fileSufix := string([]byte(filename)[fixIndex:])
			media := getMedia(fileSufix)
			fi, err := os.Open(filename)
			if err != nil {
				return "", err
			}
			io.Copy(writer, fi)
			return media, nil
		}
	}
	return "", fmt.Errorf("file not exist!")
}

func getMedia(fileFix string) string {
	media := mime.TypeByExtension(fileFix)
	if media == "" {

	}
	return media
}

//请求映射
//格式:mapper:"name(名称);url(地址);method(GET|POST);style(XML|JOSN|HTML)"。
const (
	_RequestMapperTag = "mapper"
)

func buildRequestMapperByHandlefunc(name string, url []string, input interface{}, handle HandleFunction, methods []HttpMethodType, style StyleType) *RequestMapper {
	r := &RequestMapper{}
	r.Name = name
	r.Url = url
	r.Request = input
	r.ResponseStyle = style
	r.Response = nil //输出转换器
	r.Methods = methods
	r.Handle = handle
	return r
}
func buildRequestMapperByStructTag(obj interface{}) []*RequestMapper {
	mc, ok := obj.(MutiController)
	if !ok {
		return nil
	}
	handleMap := mc.GetHandles()
	objT, objV, err := reflect2.GetStructTypeValue(obj)
	if err != nil {
		return nil
	}
	var mappers []*RequestMapper
	for i := 0; i < objT.NumField(); i++ {
		f := objT.Field(i)
		vf := objV.Field(i)
		if !vf.CanInterface() {
			continue
		}

		tag := f.Tag.Get(_RequestMapperTag)
		if len(tag) == 0 {
			continue
		} else {
			rules := strings.Split(tag, ";")
			if len(rules) > 0 {
				mapper := &RequestMapper{}
				for _, rule := range rules {
					setRequestMapper(rule, mapper, handleMap)
				}
				mappers = append(mappers, mapper)
			}

		}

	}
	return mappers
}

func setRequestMapper(exp string, mapper *RequestMapper, handles map[string]CMap) {
	vexp := strings.TrimSpace(exp)
	ruleStart := strings.Index(vexp, "(")
	var propertyName, rule string
	if ruleStart < 0 {
		propertyName = strings.ToLower(vexp)
		rule = ""
	} else {
		propertyName = strings.TrimSpace(vexp[:ruleStart])
		propertyName = strings.ToLower(propertyName)
		ruleEnd := strings.Index(vexp, ")")
		if ruleEnd < 0 {
			ruleEnd = len(vexp)
		}
		rule = strings.TrimSpace(vexp[ruleStart+1 : ruleEnd])
	}

	switch propertyName {
	case "name":
		mapper.Name = rule
		h := handles[rule]
		if h.Handle != nil {
			mapper.Handle = h.Handle
			mapper.Request = h.Parameter
		}

	case "url":
		mapper.Url = append(mapper.Url, rule)
	case "style":
		mapper.ResponseStyle = ParseHttpStyleType(rule)
		mapper.Response = convertors[mapper.ResponseStyle]
	case "method":
		mapper.Methods = append(mapper.Methods, ParseHttpMethodType(rule))

	}

}

// 结果转换处理
var convertors map[StyleType]Convertor

func init() {
	convertors = make(map[StyleType]Convertor)
	convertors[Json] = convertToJson
	convertors[Xml] = convertToXml
}

// 转json
func convertToJson(writer io.Writer, obj interface{}) error {
	if obj != nil {
		data, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		writer.Write(data)
	}
	return nil
}

// 转xml
func convertToXml(writer io.Writer, obj interface{}) error {
	if obj != nil {
		data, err := xml.Marshal(obj)
		if err != nil {
			return err
		}
		writer.Write(data)
	}
	return nil
}
