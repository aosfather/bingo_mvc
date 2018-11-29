package bingo_mvc

import (
	utils "github.com/aosfather/bingo_utils"
	"strings"
)

type AbstractDispatcher struct {
	router       routerMapper
	Port         int
	staticRoot   string
	templateRoot string
	Logger       utils.Log
}

func (this *AbstractDispatcher) SetRoot(static, template string) {
	this.staticRoot = static
	this.templateRoot = template
	this.router.SetStaticControl(this.staticRoot, this.Logger)
}

func (this *AbstractDispatcher) AddHandler(url string, handler HttpMethodHandler) {
	var rule RouterRule
	rule.Init(url, handler)
	this.router.AddRouter(&rule)
}

func (this *AbstractDispatcher) AddController(c HttpController) {
	if c != nil {
		c.Init()
		url := c.GetUrl()
		if url == "" {
			url = "/" + utils.GetRealType(c).Name()
		}
		this.AddHandler(url, c.(HttpMethodHandler))
	}
}

func (this *AbstractDispatcher) AddInterceptor(h CustomHandlerInterceptor) {

}

func (this *AbstractDispatcher) MatchURI(uri string) (*RouterRule, Params) {
	return this.router.match(uri)
}

/*
  路由映射
*/

type routerMapper struct {
	routerTree    *node
	staticHandler HttpMethodHandler
	defaultRule   *RouterRule
}

func (this *routerMapper) AddRouter(rule *RouterRule) {
	if this.routerTree == nil {
		this.routerTree = &node{}
	}
	if rule != nil {
		this.routerTree.addRoute(rule.url, rule)
	}

}

func (this *routerMapper) match(uri string) (*RouterRule, Params) {
	paramIndex := strings.Index(uri, "?")
	realuri := uri
	if paramIndex != -1 {
		realuri = strings.TrimSpace((uri[:paramIndex]))
	}

	h, p, _ := this.routerTree.getValue(realuri)
	if h == nil {
		return &RouterRule{realuri, nil, this.staticHandler}, p
	}
	return h.(*RouterRule), p
}

func (this *routerMapper) SetStaticControl(path string, l utils.Log) {
	this.staticHandler = &staticController{staticDir: path, log: l}
}
