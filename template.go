package bingo_mvc

import (
	utils "github.com/aosfather/bingo_utils"
	"html/template"
	"io"
)

/*
 模板的实现
 实现特性
  1、片段库的定义
  2、进行缓存常用的模板对象进行加速
  3、监视模板文件的变化，当模板文件变化后，摧毁缓存的模板对象
*/

type TemplateEngine struct {
	RootPath                string                        //模板根路径
	SubTemplatePath         string                        //子模板及片段定义的目录
	Suffix                  string                        //模板文件后缀
	CacheSize               int                           //缓存模板个数
	ErrorTemplate           string                        //错误模板
	useDefaultErrorTemplate bool                          //是否使用默认错误模板
	templates               map[string]*template.Template //模板缓存
}

func (this *TemplateEngine) Init() {
	this.templates = make(map[string]*template.Template)
}
func (this *TemplateEngine) Render(w io.Writer, templateName string, data interface{}) BingoError {
	//use cache
	t := this.templates[templateName]
	var templateError BingoError = nil
	if t == nil {
		templateFile := this.getRealPath(templateName)
		var err error
		if utils.IsFileExist(templateFile) {
			t, err = template.New(templateName).ParseFiles(templateFile)
			if err != nil {
				templateError = CreateError(500, "template load error:"+err.Error())
			} else {
				this.templates[templateName] = t
			}
		}
	}
	//执行template
	if t != nil {
		err := t.Execute(w, data)
		if err != nil {
			templateError = CreateErrorF(500, "template render error:%s", err.Error())
		}
	}
	return templateError
}

func (this *TemplateEngine) WriteError(w io.Writer, err BingoError) {
	if this.useDefaultErrorTemplate {
		tmpl, _ := template.New("error").Parse("<html><body><h1>{{.Code}}</h1><h3>{{.Error}}</h3></body></html>")
		tmpl.Execute(w, err)
	} else {
		e := this.Render(w, this.ErrorTemplate, err)
		if e != nil {
			this.useDefaultErrorTemplate = true
			this.WriteError(w, err)
		}
	}

}

func (this *TemplateEngine) getRealPath(templateName string) string {
	if this == nil {
		return templateName
	}

	return this.RootPath + "/" + templateName
}
