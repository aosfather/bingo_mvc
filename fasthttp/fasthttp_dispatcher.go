package fasthttp

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aosfather/bingo_mvc"
	reflect2 "github.com/aosfather/bingo_utils/reflect"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
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
		//设置服务器名称
		ctx.Response.Header.Set("Server", "bingo mvc")
		return
	}
	//获取requestmapper定义
	request := this.MatchUrl(url)
	if request == nil {
		meta, err := this.ProcessStaticUrl(url, ctx.Response.BodyWriter())
		if err != nil {
			ctx.Response.Header.Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
			ctx.Response.SetBodyString("<b>the url not found!</b>")
			ctx.Response.SetStatusCode(404)
			log.Printf("the url %s not found\n", url)
		} else {
			ctx.Response.Header.Set(bingo_mvc.CONTENT_TYPE, meta)
		}

	} else {
		if request.IsSupportMethod(bingo_mvc.ParseHttpMethodType(string(ctx.Method()))) {
			this.call(request, ctx)
		} else {
			//不支持的 http method 处理
			ctx.Response.Header.Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
			ctx.Response.SetBodyString("<b>the method not support !</b>")
			ctx.Response.SetStatusCode(405)
		}
	}

	//设置服务器名称
	ctx.Response.Header.Set("Server", "bingo mvc")

}

func (this *FastHTTPDispatcher) call(api bingo_mvc.Controller, ctx *fasthttp.RequestCtx) {
	//校验参数
	contentType := string(ctx.Request.Header.ContentType())
	//var input map[string]interface{} = make(map[string]interface{})
	//不支持文件流
	if ctx.Request.IsBodyStream() {
		ctx.Response.SetBodyString("<b>not surpport stream!</b>")
		ctx.Response.SetStatusCode(400)
		return
	}

	inputfunc := func(input interface{}) error {
		//json格式请求处理
		if strings.Contains(contentType, "application/json") {
			body := ctx.Request.Body()
			err := json.Unmarshal(body, &input)
			if err != nil {
				return err
			}
			//xml格式请求处理
		} else if strings.Contains(contentType, "text/xml") {
			body := ctx.Request.Body()
			err := xml.Unmarshal(body, &input)
			if err != nil {
				return err
			}
			//form方式
		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			this.fillByArgs(ctx.Request.PostArgs(), input)
		} else if strings.Contains(contentType, "multipart/form-data") {
			mforms, err := ctx.Request.MultipartForm()
			if err != nil {
				return err
			}

			//填充其他参数
			if len(mforms.Value) > 0 {
				inputmap := make(map[string]interface{})
				for k, v := range mforms.Value {
					inputmap[k] = strings.Join(v, "")
				}
				reflect2.FillStruct(inputmap, input)
			}
			//处理文件内容
			if fcontainer, ok := input.(bingo_mvc.FileContainer); ok {
				fs := mforms.File["file"]
				if fs != nil {
					for _, f := range fs {
						fm := &bingo_mvc.FileForm{FileName: f.Filename, FileSize: f.Size}
						content, err := f.Open()
						if err != nil {
							fm.IsError = true
							fm.Error = err.Error()
						} else {
							fm.File = content
						}
						fcontainer.AddFileForm(fm)
					}
				}
			}

		} else {
			args := ctx.QueryArgs()
			if args != nil {
				this.fillByArgs(args, input)
			}
		}
		return nil
	}

	st := this.ExecuteRequest(api, ctx.Response.BodyWriter(), &HttpContextImp{ctx}, inputfunc)
	ctx.Response.Header.Set(bingo_mvc.CONTENT_TYPE, st.GetContentType())
}

func (this *FastHTTPDispatcher) fillByArgs(args *fasthttp.Args, input interface{}) {
	if args == nil || input == nil {
		return
	}
	t := reflect2.GetRealType(input)
	var inputmap map[string]interface{}
	if t.Kind() == reflect.Map {
		inputmap = input.(map[string]interface{})
	} else {
		inputmap = make(map[string]interface{})
	}
	args.VisitAll(func(key, value []byte) {
		inputmap[string(key)] = string(value)
	})
	//如果传入的不是map类型则填充值到struct上
	if &inputmap != input {
		reflect2.FillStruct(inputmap, input)
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
