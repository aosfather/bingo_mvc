package context

import (
	container "github.com/aosfather/bingo_utils/reflect"
)

//Bean自动初始化接口
type BeanAutoInit interface {
	AutoInit() error
}

type Config interface {
	//获取属性
	GetProperty(key string) string
	//获取自定义属性
	GetPropertyForCustom(key string) string
}

//初始化函数
type InitProcessFunction func(f container.StoreFunction) (string, interface{})

type ApplicationContext struct {
	initfunctions []InitProcessFunction
	config        Config
	services      container.InjectMan
	holder        container.ValuesHolder
}

func (this *ApplicationContext) AddProcessFunction(p ...InitProcessFunction) {
	this.initfunctions = append(this.initfunctions, p...)
}

func (this *ApplicationContext) shutdown() {

	//关闭所有service

}

func (this *ApplicationContext) RegisterService(name string, service interface{}) {
	if name != "" && service != nil {
		instance := this.services.GetObjectByName(name)
		if instance == nil {
			this.holder.ProcessValueTag(service)
			this.services.AddObjectByName(name, service)
		}
	}

}

func (this *ApplicationContext) GetService(name string) interface{} {
	if name != "" {
		return this.services.GetObjectByName(name)
	}
	return nil
}

func (this *ApplicationContext) init(config Config) {
	this.config = config
	this.services.Init(nil)
	this.services.AddObject(this)
	this.holder.InitByFunction(this.config.GetPropertyForCustom)
	//初始化function
	this.initByProcessFunctions()

}

//轮询设置的初始化函数，如果存在则进行初始化，并加入到service中
func (this *ApplicationContext) initByProcessFunctions() {
	if this.initfunctions != nil && len(this.initfunctions) > 0 {
		for _, initfun := range this.initfunctions {
			name, bean := initfun(this.config.GetProperty)
			if bean == nil {
				continue
			}
			if name == "" {
				this.services.AddObject(bean)
			} else {
				this.services.AddObjectByName(name, bean)
			}
		}
	}
}

//结束加载
func (this *ApplicationContext) Finish() {
	fun := func(a interface{}) {
		this.holder.ProcessValueTag(a)
	}
	this.services.Inject(fun)
}

//组装对象
func (this *ApplicationContext) Inject(c interface{}) error {
	this.holder.ProcessValueTag(c)
	this.services.InjectObject(c)
	auto, ok := c.(BeanAutoInit)
	if ok {
		err := auto.AutoInit()
		return err
	}
	return nil
}
