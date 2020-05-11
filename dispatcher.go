package bingo_mvc

import (
	utils "github.com/aosfather/bingo_utils"
	"io"
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
func (this *AbstractDispatcher) ExecuteRequest(r *RequestMapper, writer io.Writer, request func(key string) interface{}, input func(interface{}) error) {
	if !this.preHandle(writer, request) {
		return
	}

	//执行handler处理
	r.Select(writer, input)
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
