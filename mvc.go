package bingo_mvc

import (
	utils "github.com/aosfather/bingo_utils"
	"net/http"
)

type RouterRule struct {
	url           string
	convert       *ResponseConverter
	methodHandler HttpMethodHandler
}

func (this *RouterRule) Init(url string, handle HttpMethodHandler) {
	this.url = url
	this.methodHandler = handle

}

//type BeanFactory interface {
//	GetLog(module string) utils.Log
//	GetService(name string) interface{}
//	GetSession() *sql.TxSession
//}

type Context interface {
	GetCookie(key string) string
}

//返回结果转换器，用于输出返回结果
type ResponseConverter interface {
	Convert(writer http.ResponseWriter, obj interface{})
}

type HttpMethodHandler interface {
	GetSelf() interface{}
	GetParameType(method string) interface{}
	Get(c Context, p interface{}) (interface{}, BingoError)
	Post(c Context, p interface{}) (interface{}, BingoError)
	Put(c Context, p interface{}) (interface{}, BingoError)
	Delete(c Context, p interface{}) (interface{}, BingoError)
}

type HandlerInterceptor interface {
	PreHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule) bool
	PostHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule, mv *ModelView) BingoError
	AfterCompletion(writer http.ResponseWriter, request *http.Request, handler *RouterRule, err BingoError) BingoError
}

type HttpController interface {
	Init()
	GetUrl() string
	//SetBeanFactory(f BeanFactory)
}

type Controller struct {
	//factory BeanFactory
}

func (this *Controller) Init() {

}
func (this *Controller) GetUrl() string {
	return ""
}

//func (this *Controller) SetBeanFactory(f BeanFactory) {
//	this.factory = f
//}
//func (this *Controller) GetBeanFactory() BeanFactory {
//	return this.factory
//}
func (this *Controller) GetSelf() interface{} {
	return this
}

func (this *Controller) GetParameType(method string) interface{} {
	return this

}
func (this *Controller) Get(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(Code_NOT_ALLOWED, "method not allowed!")

}
func (this *Controller) Post(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(Code_NOT_ALLOWED, "method not allowed!")
}
func (this *Controller) Put(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(Code_NOT_ALLOWED, "method not allowed!")
}
func (this *Controller) Delete(c Context, p interface{}) (interface{}, BingoError) {
	return nil, utils.CreateError(Code_NOT_ALLOWED, "method not allowed!")
}

type SimpleController struct {
	Controller
}

func (this *SimpleController) Post(c Context, p interface{}) (interface{}, BingoError) {
	return this.Get(c, p)
}
func (this *SimpleController) Put(c Context, p interface{}) (interface{}, BingoError) {
	return this.Get(c, p)
}
func (this *SimpleController) Delete(c Context, p interface{}) (interface{}, BingoError) {
	return this.Get(c, p)
}
