package http

import (
	"github.com/aosfather/bingo_mvc"
	"net/http"
)

type HttpContextImp struct {
	request  *http.Request
	response http.ResponseWriter
}

func (this *HttpContextImp) RequestHeaderRead(key string) string {
	return this.request.Header.Get(key)

}
func (this *HttpContextImp) ResponseHeaderwrite(key string, v string) error {
	this.response.Header().Set(key, v)
	return nil
}

func (this *HttpContextImp) CookieRead(key string) map[bingo_mvc.CookieKey]interface{} {
	cookie, err := this.request.Cookie(key)
	if err != nil {
		return nil
	}
	v := make(map[bingo_mvc.CookieKey]interface{})
	v[bingo_mvc.CK_Name] = cookie.Name
	v[bingo_mvc.CK_Value] = cookie.Value
	v[bingo_mvc.CK_Path] = cookie.Path
	v[bingo_mvc.CK_MaxAge] = cookie.MaxAge
	v[bingo_mvc.CK_HttpOnly] = cookie.HttpOnly
	return v
}
func (this *HttpContextImp) CookieWrite(key string, value map[bingo_mvc.CookieKey]interface{}) error {
	c := http.Cookie{}
	c.Name = value[bingo_mvc.CK_Name].(string)
	c.Value = value[bingo_mvc.CK_Value].(string)
	c.Path = value[bingo_mvc.CK_Path].(string)
	c.MaxAge = value[bingo_mvc.CK_MaxAge].(int)
	c.HttpOnly = value[bingo_mvc.CK_HttpOnly].(bool)
	http.SetCookie(this.response, &c)
	return nil
}
