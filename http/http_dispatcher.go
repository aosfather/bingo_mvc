package http

import (
	"github.com/aosfather/bingo_mvc"
	"io/ioutil"
	"net/http"
	"strconv"
)

type HttpDispatcher struct {
	bingo_mvc.AbstractDispatcher
	server *http.Server
}

func (this *HttpDispatcher) Run() {
	if this.server != nil {
		return
	}
	if this.Port <= 0 {
		this.Port = bingo_mvc.Default_Port
	}

	this.server = &http.Server{Addr: ":" + strconv.Itoa(this.Port), Handler: this}
	this.server.ListenAndServe()
}

func (this *HttpDispatcher) shutdown() {
	if this.server != nil {
		this.server.Shutdown(nil)
	}
}

func (this *HttpDispatcher) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	url := request.RequestURI
	if url == "/favicon.ico" {
		ico, _ := ioutil.ReadFile("favicon.ico")
		writer.Write(ico)
		writer.Header().Set(bingo_mvc.CONTENT_TYPE, "image/x-icon")
	}
}
