package bingo_mvc

import (
	"fmt"
	log "github.com/aosfather/bingo_utils"
	"io"
	"runtime"
)

//协议头接口
type HttpHeaderFace interface {
	RequestHeaderRead(key string) string
	GetRequestURI() string
	ResponseHeaderwrite(key string, v string) error
}

type HttpContext interface {
	HttpHeaderFace
	CookieFace
}

type Interceptor interface {
	PreHandle(writer io.Writer, context HttpContext) bool
	InputProcess(context HttpContext, input interface{}) error
	PostHandle(writer io.Writer, context HttpContext, mv *ModelView) BingoError
	AfterCompletion(writer io.Writer, context HttpContext, err BingoError) BingoError
}

type AbstractDispatcher struct {
	Port            int
	dispatchManager *DispatchManager
	interceptors    []Interceptor
	templateManager *TemplateEngine
	static          *staticControl
	sessionManager  *SessionManager
}

func (this *AbstractDispatcher) HandlePainc(after func(v interface{})) {
	if err := recover(); err != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		log.Err("http: panic serving :", err, " ,statck:\n", string(buf))
		if after != nil {
			after(err)
		}
	}
}

func (this *AbstractDispatcher) AddInterceptor(ins ...Interceptor) {
	if ins != nil && len(ins) > 0 {
		this.interceptors = append(this.interceptors, ins...)
	}
}
func (this *AbstractDispatcher) ConfigPort(p int) {
	this.Port = p
}

func (this *AbstractDispatcher) ConfigStatic(root string) {
	if root != "" {
		this.static = &staticControl{root}
	}
}
func (this *AbstractDispatcher) ConfigTemplate(root string, suffix string) {
	this.templateManager = &TemplateEngine{}
	this.templateManager.RootPath = root
	this.templateManager.Suffix = suffix
	this.templateManager.Init()
}
func (this *AbstractDispatcher) SetSessionManager(s *SessionManager) {
	this.sessionManager = s
}
func (this *AbstractDispatcher) SetDispatchManager(d *DispatchManager) {
	this.dispatchManager = d
}

func (this *AbstractDispatcher) AddRequestMapper(r *RequestMapper) {
	if r == nil {
		return
	}

	if this.dispatchManager == nil {
		this.dispatchManager = &DispatchManager{}
		this.dispatchManager.Init()
	}

	log.Debug(r.ResponseStyle)
	//使用模板来默认处理html格式
	if r.ResponseStyle == UrlForm {
		r.Response = this.convertToHtmlByTemplate
	}

	this.dispatchManager.AddRequestMapper("", r)
}
func (this *AbstractDispatcher) AddController(domain string, name string, url string, control Controller) {
	if control == nil {
		return
	}

	if this.dispatchManager == nil {
		this.dispatchManager = &DispatchManager{}
		this.dispatchManager.Init()
	}
	this.dispatchManager.AddApi(domain, name, url, control)
}
func (this *AbstractDispatcher) ProcessStaticUrl(url string, writer io.Writer) (string, error) {
	if this.static != nil {
		return this.static.Getstaticfile(url, writer)
	}
	return "", nil

}

func (this *AbstractDispatcher) MatchUrl(u string) Controller {
	if this.dispatchManager != nil {
		return this.dispatchManager.GetController("", u)
	}
	return nil
}

//执行请求
func (this *AbstractDispatcher) ExecuteRequest(r Controller, writer io.Writer, context HttpContext, input func(interface{}) error) StyleType {
	if !this.preHandle(writer, context) {
		return UrlForm
	}
	process := func(in interface{}) error {
		err := input(in)
		if err != nil {
			return err
		}
		return this.inputProcess(context, in)
	}

	st := r.Select(writer, process)
	err := this.postHandle(writer, context, nil)

	//请求处理完后拦截器进行处理
	this.afterCompletion(writer, context, err)
	return st

}
func (this *AbstractDispatcher) inputProcess(context HttpContext, input interface{}) error {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			err := h.InputProcess(context, input)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *AbstractDispatcher) preHandle(writer io.Writer, context HttpContext) bool {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			if !h.PreHandle(writer, context) {
				return false
			}
		}
	}
	return true
}
func (this *AbstractDispatcher) postHandle(writer io.Writer, context HttpContext, mv *ModelView) BingoError {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			err := h.PostHandle(writer, context, mv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (this *AbstractDispatcher) afterCompletion(writer io.Writer, context HttpContext, err BingoError) BingoError {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			e := h.AfterCompletion(writer, context, err)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (this *AbstractDispatcher) AddRequestMapperByHandleFunction(name string, url []string, input interface{}, handle HandleFunction, methods []HttpMethodType) {
	r := buildRequestMapperByHandlefunc(name, url, input, handle, methods, UrlForm)
	this.AddRequestMapper(r)
}

/**
  通过mapper 的struct tag标签加入映射
*/
func (this *AbstractDispatcher) AddRequestMapperBystruct(target interface{}) {
	mappers := buildRequestMapperByStructTag(target)

	if mappers != nil && len(mappers) > 0 {
		for _, mapper := range mappers {
			this.AddRequestMapper(mapper)
		}
	}

}

//使用模板引擎进行转换
func (this *AbstractDispatcher) convertToHtmlByTemplate(writer io.Writer, obj interface{}) error {
	view, ok := obj.(ModelView)
	if ok {
		name := view.View
		if this.templateManager != nil && name != "" {
			err := this.templateManager.Render(writer, name, view.Model)
			if err != nil {
				return fmt.Errorf("code:%d msg:%s", err.Code(), err.Error())
			}
			return nil

		}

	}

	//不需要使用模板，或者模板引擎为空
	text, ok := obj.(string)
	if ok {
		writer.Write([]byte(text))
	} else {
		writer.Write([]byte(fmt.Sprintf("%v", obj)))
	}

	return nil
}
