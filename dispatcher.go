package bingo_mvc

import (
	utils "github.com/aosfather/bingo_utils"
	"io"
	"strings"
)

type Interceptor interface {
	PreHandle(writer io.Writer, request func(key string) interface{}) bool
	PostHandle(writer io.Writer, request func(key string) interface{}, mv *ModelView) BingoError
	AfterCompletion(writer io.Writer, request func(key string) interface{}, err BingoError) BingoError
}

type AbstractDispatcher struct {
	Port            int
	Logger          utils.Log
	dispatchManager *DispatchManager
	interceptors    []Interceptor
}

func (this *AbstractDispatcher) SetDispatchManager(d *DispatchManager) {
	this.dispatchManager = d
}

func (this *AbstractDispatcher) AddRequestMapper(r *RequestMapper) {
	if this.dispatchManager == nil {
		this.dispatchManager = &DispatchManager{}
	}
	this.dispatchManager.AddRequestMapper("", r)
}

func (this *AbstractDispatcher) MatchUrl(u string) *RequestMapper {
	if this.dispatchManager != nil {
		return this.dispatchManager.GetRequestMapper("", u)
	}
	return nil
}

//执行请求
func (this *AbstractDispatcher) ExecuteRequest(r *RequestMapper, writer io.Writer, request func(key string) interface{}) {
	if !this.preHandle(writer, request) {
		return
	}

	//执行handler处理
	r.Select(writer, nil)
	err := this.postHandle(writer, request, nil)

	//请求处理完后拦截器进行处理
	this.afterCompletion(writer, request, err)

}

func (this *AbstractDispatcher) preHandle(writer io.Writer, request func(key string) interface{}) bool {
	if this.interceptors != nil && len(this.interceptors) > 0 {
		for _, h := range this.interceptors {
			if !h.PreHandle(writer, request) {
				return false
			}
		}
	}
	return true
}
func (this *AbstractDispatcher) postHandle(writer io.Writer, request func(key string) interface{}, mv *ModelView) BingoError {
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
func (this *AbstractDispatcher) afterCompletion(writer io.Writer, request func(key string) interface{}, err BingoError) BingoError {
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
  路由映射
*/

type routerMapper struct {
	routerTree    *node
	staticHandler HttpMethodHandler
	defaultRule   *RouterRule
}

func (this *routerMapper) AddRouter(rule *RouterRule) {
	if this.routerTree == nil {
		this.routerTree = &node{}
	}
	if rule != nil {
		this.routerTree.addRoute(rule.url, rule)
	}

}

func (this *routerMapper) match(uri string) (*RouterRule, Params) {
	paramIndex := strings.Index(uri, "?")
	realuri := uri
	if paramIndex != -1 {
		realuri = strings.TrimSpace((uri[:paramIndex]))
	}

	h, p, _ := this.routerTree.getValue(realuri)
	if h == nil {
		return &RouterRule{realuri, nil, this.staticHandler}, p
	}
	return h.(*RouterRule), p
}

func (this *routerMapper) SetStaticControl(path string, l utils.Log) {
	this.staticHandler = &staticController{staticDir: path, log: l}
}
