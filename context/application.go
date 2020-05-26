package context

import (
	"fmt"
	"github.com/aosfather/bingo_utils/files"
	container "github.com/aosfather/bingo_utils/reflect"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//Bean自动初始化接口
type BeanAutoInit interface {
	AutoInit() error
}

//初始化函数
type InitProcessFunction func(f container.StoreFunction) (string, interface{})

type ApplicationContext struct {
	initfunctions []InitProcessFunction
	config        map[interface{}]interface{}
	services      container.InjectMan
	holder        container.ValuesHolder
}

func (this *ApplicationContext) AddProcessFunction(p ...InitProcessFunction) {
	this.initfunctions = append(this.initfunctions, p...)
}

func (this *ApplicationContext) shutdown() {

	//关闭所有service

}

//不能获取bingo自身的属性，只能获取应用自身的扩展属性
func (this *ApplicationContext) GetPropertyFromConfig(key string) string {
	if strings.HasPrefix(key, "bingo.") {
		return ""
	}
	return this.getProperty(key)
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

func (this *ApplicationContext) getProperty(key string) string {
	if this.config == nil {
		return ""
	}
	v, ok := this.config[key]
	if ok {
		return v.(string)
	} else {
		if strings.Index(key, ".") > 0 {
			keys := strings.Split(key, ".")
			return this.getvalue(this.config, keys, 0)
		}
	}
	return ""
}

func (this *ApplicationContext) getvalue(m map[interface{}]interface{}, keys []string, index int) string {
	v, ok := m[keys[index]]
	if ok {
		if value, ok := v.(string); ok {
			return value
		}

		if value, ok := v.(map[interface{}]interface{}); ok {
			return this.getvalue(value, keys, index+1)
		}

		return fmt.Sprintf("%v", v)

	}
	return ""
}
func (this *ApplicationContext) init(file string) {
	if file != "" && files.IsFileExist(file) {
		f, err := os.Open(file)
		if err == nil {
			txt, _ := ioutil.ReadAll(f)
			err := yaml.Unmarshal(txt, &this.config)
			if err != nil {
				log.Println(err.Error())
				panic("load config file error!")
			}
		}

	}
	if this.config == nil {
		this.config = make(map[interface{}]interface{})
	}

	this.services.Init(nil)
	this.services.AddObject(this)
	this.holder.InitByFunction(this.GetPropertyFromConfig)
	//初始化function
	this.initByProcessFunctions()

}

//轮询设置的初始化函数，如果存在则进行初始化，并加入到service中
func (this *ApplicationContext) initByProcessFunctions() {
	if this.initfunctions != nil && len(this.initfunctions) > 0 {
		for _, initfun := range this.initfunctions {
			name, bean := initfun(this.getProperty)
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
