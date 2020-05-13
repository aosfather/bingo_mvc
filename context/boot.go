package context

import (
	"github.com/aosfather/bingo_mvc"
	"github.com/aosfather/bingo_utils/files"
	logs "github.com/aosfather/bingo_utils/log"
	"github.com/aosfather/bingo_utils/reflect"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

/**
  启动包装类，用于简单启动系统
*/
//load函数，如果加载成功返回true，否则返回FALSE
type Dispatch interface {
	ConfigPort(p int)
	ConfigStatic(root string)
	ConfigTemplate(root string, suffix string)
	AddRequestMapperBystruct(target interface{})
	AddRequestMapperByHandleFunction(name string, url []string, input interface{}, handle bingo_mvc.HandleFunction, methods []bingo_mvc.HttpMethodType)
	Run()
}
type OnBootLoad func() []interface{}
type OnDestoryHandler func() bool //shutdown的handler，用于处理关闭服务的自定义动作
type Boot struct {
	applicationContext *ApplicationContext
	dispatch           Dispatch
	onShutdown         OnDestoryHandler
	onloads            []OnBootLoad
}

func (this *Boot) Init(d Dispatch, lfunc ...OnBootLoad) {
	this.dispatch = d
	this.onloads = append(this.onloads, lfunc...)
}

func (this *Boot) Start() {
	var configfile string
	if len(os.Args) > 1 {
		configfile = os.Args[1]
	} else {
		configfile = "bingo.yaml"
	}

	if !files.IsFileExist(configfile) {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		configfile = dir + "/" + configfile
	}

	this.StartByConfigFile(configfile)
}

func (this *Boot) StartByConfigFile(filename string) {
	go this.signalListen()
	this.applicationContext = &ApplicationContext{}
	//加入数据库处理模块
	this.applicationContext.AddProcessFunction(InitDatasource, this.initDispatch)

	this.applicationContext.init(filename)
	//加载factory
	if this.onloads != nil && len(this.onloads) > 0 {
		for _, onload := range this.onloads {
			objs := onload()
			if len(objs) > 0 {
				for _, obj := range objs {
					this.applicationContext.services.AddObject(obj)
					_, ok := obj.(bingo_mvc.MutiController)
					if ok {
						this.dispatch.AddRequestMapperBystruct(obj)
					}
				}

			}

		}
	}

	this.applicationContext.Finish()
	if this.dispatch != nil {
		this.dispatch.Run()
	}

}

func (this *Boot) initDispatch(f reflect.StoreFunction) (string, interface{}) {
	if this.dispatch != nil {
		port, _ := strconv.Atoi(f("bingo.port"))
		log.Println(port)
		this.dispatch.ConfigPort(port)
		this.dispatch.ConfigStatic(f("bingo.static"))
		this.dispatch.ConfigTemplate(f("bingo.template"), f("bingo.template_fix"))

	}
	return "", nil
}

//初始化日志工厂
func (this *Boot) initLogFactory(f reflect.StoreFunction) (string, interface{}) {
	logfactory := &logs.LogFactory{}
	logfactory.SetConfig(logs.LogConfig{true, f("bingo.log.file")})
	return "", logfactory
}

const (
	SIGUSR1 = syscall.Signal(0x1e)
	SIGUSR2 = syscall.Signal(0x1f)
)

//监听被kill的信号，当被kill的时候执行处理
func (this *Boot) signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, SIGUSR1, SIGUSR2)
	s := <-c
	//收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
	log.Printf("get signal:%s", s)
	this.processShutdown()

}

func (this *Boot) processShutdown() {
	//处理关闭操作
	//关闭service
	this.applicationContext.shutdown()
	//处理自定义关闭操作
	if this.onShutdown != nil {
		this.onShutdown()
	}

	os.Exit(0)

}
