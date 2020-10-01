package bingo_mvc

import (
	log "github.com/aosfather/bingo_utils"
	. "strings"
)

/**
  请求分发管理
*/
type DispatchManager struct {
	domainNode  map[string]*node  //特定域名下的node
	defaultNode *node             //默认
	apiMap      map[string]Handle //api列表
}

func (this *DispatchManager) Init() {
	this.domainNode = make(map[string]*node)
	this.apiMap = make(map[string]Handle)
	this.defaultNode = &node{}
}

//根据域名和url获取对应的API
func (this *DispatchManager) GetApi(domain, url string) (Handle,Params) {
	node := this.domainNode[domain]
	if node == nil {
		node = this.defaultNode
	}

	if node != nil {
		paramIndex := Index(url, "?")
		realuri := url
		if paramIndex != -1 {
			realuri = TrimSpace((url[:paramIndex]))
		}

		h, p, _ := node.getValue(realuri)
		if h != nil {
			key := h.(string)
			return this.apiMap[key],p
		}

	}
	return nil,nil
}

/**
增加单个api的映射
*/
func (this *DispatchManager) AddApi(domain string, name, url string, handle Handle) {
	if handle == nil || name == "" || url == "" {
		return
	}

	this.apiMap[name] = handle
	var apiNode *node
	if domain == "" {
		apiNode = this.defaultNode
	} else {
		//处理不同的域名的映射
		if apiNode == nil {
			apiNode = this.domainNode[domain]
			if apiNode == nil {
				apiNode = &node{}
				this.domainNode[domain] = apiNode
			}
		}

	}

	if handle != nil {
		apiNode.addRoute(url, name)
	}
}

/**
 新增一个requestmapper的映射。
一个requestmapper会对应多个url
*/
func (this *DispatchManager) AddRequestMapper(domain string, r *RequestMapper) {
	if r != nil {
		for _, url := range r.Url {
			log.Debug(url)
			this.AddApi(domain, r.Name, url, r)
		}
	}
}

func (this *DispatchManager) GetRequestMapper(domain, url string) *RequestMapper {
	r,_ := this.GetApi(domain, url)
	if r != nil {
		return r.(*RequestMapper)
	}
	return nil
}

func (this *DispatchManager) GetController(domain, url string) (Controller,Params) {
	r,p := this.GetApi(domain, url)
	if r != nil {
		return r.(Controller),p
	}
	return nil,nil
}
