package bingo_mvc

import (
	"fmt"
	utils "github.com/aosfather/bingo_utils"
	"io"
	"os"
	"reflect"
	"strings"
)

//url元信息描述
type RequestMapper struct {
	Name    string   //名称
	Url     []string //url路径
	Methods []HttpMethodType
	//请求参数类型
	Request interface{}
	//返回值处理器
	Response Convertor
	Handle   HandleFunction
}

func (this *RequestMapper) IsSupportMethod(m HttpMethodType) bool {
	for _, sm := range this.Methods {
		if sm == m {
			return true
		}
	}
	return false
}

func (this *RequestMapper) Select(writer io.Writer, input func(interface{}) error) {
	t := reflect.TypeOf(this.Request)
	if t.Kind() == reflect.Ptr { //指针类型获取真正type需要调用Elem
		t = t.Elem()
	}
	paramter := reflect.New(t).Interface()
	err := input(paramter)
	if err != nil {

	} else {
		result := this.Handle(paramter)
		err = this.Response(writer, result)
		//错误处理进行输出

	}

}

//结果集转换器
type Convertor func(writer io.Writer, obj interface{}) error

//Request响应函数
type HandleFunction func(input interface{}) interface{}

//控制器
type Control interface {
	Select(writer io.Writer, input func(interface{}) error)
	IsSupportMethod(m HttpMethodType) bool
}

//静态资源处理
type staticControl struct {
	root string
}

func (this *staticControl) Getstaticfile(url string, writer io.Writer) (string, error) {
	filename := this.root + url
	if utils.IsFileExist(filename) {
		fixIndex := strings.LastIndex(filename, ".")
		fileSufix := string([]byte(filename)[fixIndex:])
		media := getMedia(fileSufix)
		fi, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		io.Copy(writer, fi)
		return media, nil
	}
	return "", fmt.Errorf("file not exist!")
}
