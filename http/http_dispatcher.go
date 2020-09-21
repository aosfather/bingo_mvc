package http

import (
	. "bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aosfather/bingo_mvc"
	log "github.com/aosfather/bingo_utils"
	"github.com/aosfather/bingo_utils/reflect"
	"io/ioutil"
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
	defer this.HandlePainc(func(v interface{}) {
		writer.Header().Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
		writer.WriteHeader(500)
		writer.Write([]byte(fmt.Sprintf("<b>runtime error!</b><p>%v</p>", v)))
	})
	url := request.RequestURI
	//设置服务器名称
	writer.Header().Set("Server", "bingo mvc")
	if url == "/favicon.ico" {
		ico, _ := ioutil.ReadFile("favicon.ico")
		writer.Header().Set(bingo_mvc.CONTENT_TYPE, "image/x-icon")
		writer.WriteHeader(200)
		writer.Write(ico)
		return
	}
	//获取requestmapper定义
	requestMapper := this.MatchUrl(url)
	if requestMapper == nil {
		buffer := new(Buffer)
		meta, err := this.ProcessStaticUrl(url, buffer)
		if err != nil {
			writer.Header().Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
			writer.WriteHeader(404)
			writer.Write([]byte("<b>the url not found!</b>"))
			log.Debugf("the url %s not found\n", url)
		} else {
			writer.Header().Set(bingo_mvc.CONTENT_TYPE, meta)
			writer.Write(buffer.Bytes())
		}
	} else {
		if requestMapper.IsSupportMethod(bingo_mvc.ParseHttpMethodType(request.Method)) {
			this.call(requestMapper, request, writer)
		} else {
			//不支持的 http method 处理
			writer.Header().Set(bingo_mvc.CONTENT_TYPE, "text/html;charset=utf-8")
			writer.WriteHeader(405)
			writer.Write([]byte("<b>the method not support !</b>"))
		}
	}

}

func (this *HttpDispatcher) call(api bingo_mvc.Controller, request *http.Request, writer http.ResponseWriter) {
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
		} else if strings.Contains(contentType, "multipart/form-data") {
			err := request.ParseMultipartForm(256)
			if err != nil {
				return err
			}

			mforms := request.MultipartForm
			//填充其它参数
			if len(mforms.Value) > 0 {
				inputmap := make(map[string]interface{})
				for k, v := range mforms.Value {
					inputmap[k] = strings.Join(v, "")
				}
				reflect.FillStruct(inputmap, input)
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
			if request.Form == nil {
				request.ParseForm()
			}
			this.fillByForm(request.Form, input)

		}
		return nil
	}
	buffer := new(Buffer)
	st := this.ExecuteRequest(api, buffer, &HttpContextImp{request, writer}, inputfunc)
	writer.Header().Set(bingo_mvc.CONTENT_TYPE, st.GetContentType())
	writer.Write(buffer.Bytes())
}

func (this *HttpDispatcher) fillByForm(form url.Values, input interface{}) {
	if reflect.IsMap(input) {
		if sr, ok := input.(map[string]interface{}); ok {
			for key, _ := range form {
				sr[key] = form.Get(key)
			}
		}
	} else {
		reflect.FillStructByForm(form, input)
	}
}
