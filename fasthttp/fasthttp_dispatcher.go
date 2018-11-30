package fasthttp

import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
	"github.com/valyala/fasthttp"
)

/**
基于fasthttp实现
*/
type FastHTTPDispatcher struct {
	bingo_mvc.AbstractDispatcher
	server *fasthttp.Server
}

func (this *FastHTTPDispatcher) handle(ctx *fasthttp.RequestCtx) {
	ctx.Response.AppendBodyString("hello bingo")
}

func (this *FastHTTPDispatcher) Init() {
	if this.Port == 0 {
		this.Port = 8990
	}

}

func (this *FastHTTPDispatcher) Run() {
	this.server = &fasthttp.Server{Handler: this.handle}
	addr := fmt.Sprintf("0.0.0.0:%d", this.Port)
	this.server.ListenAndServe(addr)
}

func (this *FastHTTPDispatcher) Shutdown() {
	if this.server != nil {
		this.server.Shutdown()
	}
}
