package fasthttp

import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
	"github.com/valyala/fasthttp"
	"io/ioutil"
)

/**
基于fasthttp实现
*/
type FastHTTPDispatcher struct {
	bingo_mvc.AbstractDispatcher
	server *fasthttp.Server
}

func (this *FastHTTPDispatcher) handle(ctx *fasthttp.RequestCtx) {
	url := string(ctx.Request.URI().RequestURI())

	//domain := string(ctx.Request.Header.Host())
	if url == "/favicon.ico" {
		ico, _ := ioutil.ReadFile("favicon.ico")
		ctx.Response.Header.Set(bingo_mvc.CONTENT_TYPE, "image/x-icon")
		ctx.Response.SetBodyRaw(ico)
	}
}

func (this *FastHTTPDispatcher) Run() {
	this.server = &fasthttp.Server{Handler: this.handle}
	if this.Port == 0 {
		this.Port = 8990
	}
	addr := fmt.Sprintf("0.0.0.0:%d", this.Port)
	this.server.ListenAndServe(addr)
}

func (this *FastHTTPDispatcher) Shutdown() {
	if this.server != nil {
		this.server.Shutdown()
	}
}
