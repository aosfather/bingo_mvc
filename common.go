package bingo_mvc

import (
	"fmt"
	"strings"
)

//带错误码的错误接口
type BingoError interface {
	error
	Code() int
}

type Error struct {
	code int
	msg  string
}

func (this Error) Code() int {
	return this.code
}

func (this Error) Error() string {
	return this.msg
}

func CreateError(c int, text string) Error {
	return Error{c, text}
}

func CreateErrorF(c int, f string, textobj ...interface{}) Error {
	text := fmt.Sprintf(f, textobj...)
	return Error{c, text}
}

const (
	//http方法
	Method_GET    = "GET"
	Method_POST   = "POST"
	Method_PUT    = "PUT"
	Method_DELETE = "DELETE"
	Method_PATCH  = "PATCH"

	//返回码
	Code_OK             = 200
	Code_CREATED        = 201
	Code_EMPTY          = 204
	Code_NOT_MODIFIED   = 304
	Code_BAD            = 400
	Code_UNAUTHORIZED   = 401
	Code_FORBIDDEN      = 403
	Code_NOT_FOUND      = 404
	Code_CONFLICT       = 409
	Code_ERROR          = 500
	Code_NOT_IMPLEMENTS = 501
	Code_NOT_ALLOWED    = 405

	//内容类型
	_CONTENT_TYPE      = "Content-Type"
	_CONTENT_TYPE_JSON = "application/json"
	_CONTENT_JSON      = "application/json;charset=utf-8"
	_CONTENT_HTML      = "text/html"
	_CONTENT_XML       = "application/xml;charset=utf-8"

	CONTENT_TYPE = "Content-Type"
)

type HttpMethodType byte

const (
	Get  HttpMethodType = 20
	Post HttpMethodType = 21
	Put  HttpMethodType = 22
	Del  HttpMethodType = 23
	Head HttpMethodType = 24
)

const (
	Default_Port = 8080
)

func ParseHttpMethodType(method string) HttpMethodType {
	method = strings.ToUpper(method)
	switch method {
	case Method_GET:
		return Get
	case Method_POST:
		return Post
	case Method_PUT:
		return Put
	case Method_DELETE:
		return Del
	}
	return Get
}

//数据格式类型
type StyleType byte

const (
	Json    StyleType = 11
	Xml     StyleType = 12
	UrlForm StyleType = 13
)

func ParseHttpStyleType(styleName string) StyleType {
	styleName = strings.ToUpper(styleName)
	switch styleName {
	case "JSON":
		return Json
	case "XML":
		return Xml
	default:
		return UrlForm
	}
}

func (this StyleType) GetContentType() string {
	switch this {
	case Json:
		return "application/json;charset=utf-8"
	case Xml:
		return "text/xml;charset=utf-8"
	case UrlForm:
		return "text/html;charset=utf-8"
	}
	return "text/html"
}
