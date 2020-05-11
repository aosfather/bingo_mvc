package http

import (
	"encoding/json"
	"encoding/xml"
	"github.com/aosfather/bingo_mvc"
	utils "github.com/aosfather/bingo_utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
		return
	}
	//获取requestmapper定义
	requestMapper := this.MatchUrl(url)
	if requestMapper == nil {
		writer.Header().Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
		writer.Write([]byte("<b>the url not found!</b>"))
		writer.WriteHeader(404)
		log.Printf("the url %s not found\n", url)
	} else {
		if requestMapper.IsSupportMethod(bingo_mvc.ParseHttpMethodType(request.Method)) {
			this.call(requestMapper, request, writer)
		} else {
			//不支持的 http method 处理
			writer.Header().Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
			writer.Write([]byte("<b>the method not support !</b>"))
			writer.WriteHeader(405)
		}
	}
	//设置服务器名称
	writer.Header().Set("Server", "bingo mvc")

}

func (this *HttpDispatcher) call(api *bingo_mvc.RequestMapper, request *http.Request, writer http.ResponseWriter) {
	//不支持文件流
	if request.MultipartForm != nil {
		writer.Header().Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
		writer.Write([]byte("<b>not surpport stream!</b>"))
		writer.WriteHeader(400)
		return
	}

	//校验参数
	contentType := request.Header.Get(bingo_mvc.CONTENT_TYPE)
	inputfunc := func(input interface{}) error {
		//json格式请求处理
		if strings.Contains(contentType, "application/json") {
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				return err
			}
			err = json.Unmarshal(body, &input)
			if err != nil {
				return err
			}
			//xml格式请求处理
		} else if strings.Contains(contentType, "text/xml") {
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				return err
			}
			err = xml.Unmarshal(body, &input)
			if err != nil {
				return err
			}
			//form方式
		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			if request.PostForm == nil {
				request.ParseForm()
			}
			this.fillByForm(request.PostForm, input)
		} else {
			if request.Form == nil {
				request.ParseForm()
			}
			this.fillByForm(request.PostForm, input)

		}
		return nil
	}
	//获取header头的信息
	headerFunc := func(key string) interface{} {
		return request.Header.Get(key)
	}

	this.ExecuteRequest(api, writer, headerFunc, inputfunc)
	writer.Header().Set(bingo_mvc.CONTENT_TYPE, api.ResponseStyle.GetContentType())
}

func (this *HttpDispatcher) fillByForm(form url.Values, input interface{}) {
	if utils.IsMap(input) {
		if sr, ok := input.(map[string]interface{}); ok {
			for key, _ := range form {
				sr[key] = form.Get(key)
			}
		}
	} else {
		utils.FillStructByForm(form, input)
	}
}