package bingo_mvc

import (
	"io"
	"net/http"
)

const (
	URL_TAG = "Url"
)

type Context interface {
	GetCookie(key string) string
}

//简单的返回结果。用于rest api方式的返回
type SimpleResult struct {
	Action    string
	Success   bool
	ErrorCode int
	Msg       string
}

type ModelView struct {
	View  string
	Model interface{}
}
type StaticResource struct {
	Type string
	Uri  string
}

type RedirectEntity struct {
	Url     string
	Code    int
	Cookies []*http.Cookie
}
type MutiStruct interface {
	GetData() interface{}
	GetDataType() string
}

type FileHandler interface {
	io.Reader
	io.Closer
}

type StaticView struct {
	Name   string      //资源名称
	Media  string      //资源类型
	Length int         //资源长度
	Reader FileHandler //资源内容
}
