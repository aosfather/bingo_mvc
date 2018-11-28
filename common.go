package bingo_mvc

type BingoError interface {
	error
	Code() int
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
)
