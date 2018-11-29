package bingo_mvc

import (
	"encoding/json"
	"encoding/xml"
	utils "github.com/aosfather/bingo_utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type HttpDispatcher struct {
	abstractDispatcher
	server         *http.Server
	interceptor    defaultHandlerInterceptor
	defaultConvert defaultResponseConverter
}

func (this *HttpDispatcher) Init() {
	if this.port == 0 {
		this.port = 8990
	}

}

func (this *HttpDispatcher) Run() {
	if this.server != nil {
		return
	}
	this.server = &http.Server{Addr: ":" + strconv.Itoa(this.port), Handler: this}
	this.server.ListenAndServe()
}

func (this *HttpDispatcher) shutdown() {
	if this.server != nil {
		this.server.Shutdown(nil)
	}
}

func (this *HttpDispatcher) doConvert(writer http.ResponseWriter, rule *RouterRule, req *http.Request, obj interface{}) {
	if err, ok := obj.(BingoError); ok {
		writer.WriteHeader(err.Code())
		obj = ModelView{"error", err}
	}

	if rule != nil && rule.convert != nil {
		(*rule.convert).Convert(writer, obj)
	} else {

		this.defaultConvert.Convert(writer, obj, req)
	}
}

func (this *HttpDispatcher) doMethod(request *http.Request, handler HttpMethodHandler, p Params) (interface{}, BingoError) {
	method := request.Method
	param := handler.GetParameType(method)
	this.parseRequest(request, p, param)
	errors := Validate(param)
	if errors != nil && len(errors) > 0 {
		var errorText string
		for _, err := range errors {
			errorText += err.Error() + ";"
		}
		return nil, utils.CreateError(400, errorText)
	}

	var context ContextImp
	context.request = request

	var result interface{}
	var err BingoError
	switch method {
	case Method_GET:
		result, err = handler.Get(&context, param)
	case Method_POST:
		result, err = handler.Post(&context, param)
	case Method_PUT:
		result, err = handler.Put(&context, param)
	case Method_DELETE:
		result, err = handler.Delete(&context, param)
	default:
		result, err = nil, utils.CreateError(405, "method not found!")
	}

	return result, err

}

func (this *HttpDispatcher) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	uri := request.RequestURI
	rule, p := this.router.match(uri)
	handler := rule.methodHandler
	this.logger.Debug("request %s", uri)
	//handler前拦截器处理
	if !this.interceptor.PreHandle(writer, request, rule) {
		return
	}

	//执行handler处理
	obj, err := this.doMethod(request, handler, p)

	//handler处理完后，拦截器进行额外补充处理
	if mv, ok := obj.(ModelView); ok {
		this.interceptor.PostHandle(writer, request, rule, &mv)
	} else {
		this.interceptor.PostHandle(writer, request, rule, nil)
	}

	//进行结果输出
	if err != nil {
		this.doConvert(writer, rule, request, err)
	} else {
		this.doConvert(writer, rule, request, obj)
	}

	//请求处理完后拦截器进行处理
	this.interceptor.AfterCompletion(writer, request, rule, err)

}

func (this *HttpDispatcher) parseRequest(request *http.Request, p Params, target interface{}) {
	//静态资源的处理
	if sr, ok := target.(*StaticResource); ok {
		sr.Type = request.Header.Get(_CONTENT_TYPE)
		sr.Uri = request.RequestURI
		return
	}

	contentType := request.Header.Get(_CONTENT_TYPE)
	if _CONTENT_TYPE_JSON == contentType || _CONTENT_JSON == contentType || strings.Contains(contentType, _CONTENT_TYPE_JSON) { //处理为json的输入
		input, err := ioutil.ReadAll(request.Body)
		this.logger.Debug(string(input))
		defer request.Body.Close()
		if err == nil {
			if request.Form == nil {
				request.ParseForm()
				addParamsToForm(request.Form, p)
			}

			utils.FillStructByForm(request.Form, target)

			jsonTarget := target
			if sr, ok := target.(MutiStruct); ok {

				jsonTarget = sr.GetData()

			}

			err = json.Unmarshal(input, jsonTarget)
			if err != nil {
				this.logger.Error("parse request body as json error:%s", err)
			}

		} else {
			this.logger.Debug("read request body error:%s", err)
		}

	} else { //标准form的处理
		if request.Form == nil {
			request.ParseForm()
			addParamsToForm(request.Form, p)
		}

		formvalues := request.Form
		this.logger.Debug("form:%s", formvalues)

		if utils.IsMap(target) {
			if sr, ok := target.(map[string]string); ok {
				for key, _ := range formvalues {
					sr[key] = formvalues.Get(key)
				}
			}
		} else {
			utils.FillStructByForm(request.Form, target)
		}

		if sr, ok := target.(MutiStruct); ok {
			input, err := ioutil.ReadAll(request.Body)
			this.logger.Debug("input body:%s", input)
			defer request.Body.Close()
			if err == nil {
				//
				if sr.GetDataType() == "json" {
					parameters := make(map[string]interface{})
					json.Unmarshal(input, &parameters)
					utils.FillStruct(parameters, sr.GetData())
				} else if sr.GetDataType() == "xml" {
					xml.Unmarshal(input, sr.GetData())
				}

			}
		}

	}

}

type ContextImp struct {
	request *http.Request
}

func (this *ContextImp) GetCookie(key string) string {
	if this.request != nil {
		cookie, _ := this.request.Cookie(key)
		if cookie != nil {
			return cookie.Value
		}

	}
	return ""
}

type CustomHandlerInterceptor interface {
	PreHandle(writer http.ResponseWriter, request *http.Request) bool
	PostHandle(writer http.ResponseWriter, request *http.Request, mv *ModelView) BingoError
	AfterCompletion(writer http.ResponseWriter, request *http.Request, err BingoError) BingoError
}

type defaultHandlerInterceptor struct {
	interceptors []CustomHandlerInterceptor
}

func (this *defaultHandlerInterceptor) addInterceptor(interceptor CustomHandlerInterceptor) {
	if this.interceptors == nil {
		this.interceptors = []CustomHandlerInterceptor{interceptor}
	} else {
		this.interceptors = append(this.interceptors, interceptor)
	}
}

func (this *defaultHandlerInterceptor) PreHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule) bool {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			if !h.PreHandle(writer, request) {
				return false
			}
		}
	}
	return true
}
func (this *defaultHandlerInterceptor) PostHandle(writer http.ResponseWriter, request *http.Request, handler *RouterRule, mv *ModelView) BingoError {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			err := h.PostHandle(writer, request, mv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (this *defaultHandlerInterceptor) AfterCompletion(writer http.ResponseWriter, request *http.Request, handler *RouterRule, err BingoError) BingoError {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			e := h.AfterCompletion(writer, request, err)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

/*
默认返回转换器
1、根据返回类型来进行转换
2、ModelView-> 走template转换
3、其它类型->走json
4、文件流的支持？
5、xml的支持?
6、图片?
*/
