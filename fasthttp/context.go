package fasthttp

import (
	"github.com/aosfather/bingo_mvc"
	"github.com/valyala/fasthttp"
)

type HttpContextImp struct {
	ctx *fasthttp.RequestCtx
}

func (this *HttpContextImp) RequestHeaderRead(key string) string {
	return string(this.ctx.Request.Header.Peek(key))

}
func (this *HttpContextImp) ResponseHeaderwrite(key string, v string) error {
	this.ctx.Response.Header.Set(key, v)
	return nil
}

func (this *HttpContextImp) CookieRead(key string) map[bingo_mvc.CookieKey]interface{} {
	value := this.ctx.Request.Header.Cookie(key)
	c := fasthttp.Cookie{}
	c.ParseBytes(value)
	v := make(map[bingo_mvc.CookieKey]interface{})
	v[bingo_mvc.CK_Name] = c.Key()
	v[bingo_mvc.CK_Value] = c.Value()
	v[bingo_mvc.CK_Path] = c.Path()
	v[bingo_mvc.CK_MaxAge] = c.MaxAge()
	v[bingo_mvc.CK_HttpOnly] = c.HTTPOnly()
	return v
}
func (this *HttpContextImp) CookieWrite(key string, value map[bingo_mvc.CookieKey]interface{}) error {
	c := fasthttp.Cookie{}
	c.SetKey(value[bingo_mvc.CK_Name].(string))
	c.SetValue(value[bingo_mvc.CK_Value].(string))
	c.SetPath(value[bingo_mvc.CK_Path].(string))
	c.SetMaxAge(value[bingo_mvc.CK_MaxAge].(int))
	c.SetHTTPOnly(value[bingo_mvc.CK_HttpOnly].(bool))
	this.ctx.Response.Header.SetCookie(&c)
	return nil
}
