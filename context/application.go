package context

import (
	"encoding/json"
	"github.com/aosfather/bingo_utils/files"
	"github.com/aosfather/bingo_utils/log"
	container "github.com/aosfather/bingo_utils/reflect"
	"io/ioutil"
	"os"
	"strings"
)

//Bean自动初始化接口
type BeanAutoInit interface {
	AutoInit() error
}

type ApplicationContext struct {
	config     map[string]string
	logfactory *log.LogFactory
	services   container.InjectMan
	holder     container.ValuesHolder
}

func (this *ApplicationContext) shutdown() {

	//关闭所有service

	this.logfactory.Close()

}

func (this *ApplicationContext) GetLog(module string) log.Log {
	return this.logfactory.GetLog(module)
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
	return this.config[key]
}
func (this *ApplicationContext) init(file string) {
	if file != "" && files.IsFileExist(file) {
		f, err := os.Open(file)
		if err == nil {
			txt, _ := ioutil.ReadAll(f)
			json.Unmarshal(txt, &this.config)
		}

	}
	if this.config == nil {
		this.config = make(map[string]string)
	}

	this.services.Init(nil)
	this.services.AddObject(this)
	this.holder.InitByFunction(this.GetPropertyFromConfig)
	this.initLogFactory()
}

func (this *ApplicationContext) initLogFactory() {
	this.logfactory = &log.LogFactory{}
	this.logfactory.SetConfig(log.LogConfig{true, this.config["bingo.log.file"]})
}

//结束加载
func (this *ApplicationContext) Finish() {
	this.services.Inject()
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
