package flow

import (
	"fmt"
	"github.com/aosfather/bingo_utils/contain"
	"time"
)

//任务接口
type TaskHandle interface {
	SetNotifylistener(l NotifyTaskStatus)
	GetName() string
	GetTaskDefine() TaskDefine
	Execute(flowid, taskid string, parameter ...Parameter) error
}

type NotifyTaskStatus func(flowid string, taskid int, success bool, err string, parameter ...Parameter)

//管理器:任务定义管理器、流程管理器
type _TaskManager struct {
	tasks   map[string]TaskHandle  //任务实现
	defines map[string]*TaskDefine //任务定义

}

func (this *_TaskManager) Init() {
	if this.defines == nil {
		this.defines = make(map[string]*TaskDefine)
	}

	if this.tasks == nil {
		this.tasks = make(map[string]TaskHandle)
	}
}

func (this *_TaskManager) addHandle(h TaskHandle) {
	this.tasks[h.GetName()] = h
	define := h.GetTaskDefine()
	this.defines[define.Name] = &define
}

func (this *_TaskManager) getTaskDefine(name string) *TaskDefine {
	return nil
}

func (this *_TaskManager) execute(name string, flowid string, taskid int, parameter ...Parameter) error {

	return nil
}

type FlowManager struct {
	flows map[string]*Flow
}

func (this *FlowManager) Init() {
	if this.flows == nil {
		this.flows = make(map[string]*Flow)
	}
}

//发布流程
func (this *FlowManager) Publish(f *Flow) error {
	if f != nil {
		fname := f.Name
		oldflow := this.flows[fname]

		if oldflow == nil {
			this.flows[fname] = f
			return nil
		} else {
			//版本高才会更新
			if f.Version > oldflow.Version {
				this.flows[fname] = f
				//TODO 备份旧版本
				return nil
			}
		}

	}

	return fmt.Errorf("flow object is nil!")
}

func (this *FlowManager) GetFlow(name string) *Flow {
	if name != "" {
		return this.flows[name]
	}

	return nil
}

//删除流程
func (this *FlowManager) Remove(name string) error {
	if name != "" {
		delete(this.flows, name)
		return nil
	}

	return fmt.Errorf("flow name is empty!")
}

//元信息
type WorkflowMetaManager interface {
	Publish(flow *Flow) error
	GetFlow(name string) *Flow
	Remove(name string) error
}

//状态存储
type InstanceStore interface {
}

type WorkflowManager struct {
	FlowMeta    WorkflowMetaManager
	Store       InstanceStore
	taskmanager *_TaskManager
	Prefix      string
	instanceNo  int
	_instance   *contain.Cache
}

func (this *WorkflowManager) Init() {
	this.taskmanager = &_TaskManager{}
	this.taskmanager.Init()
	this._instance = contain.New(10*time.Minute, 0)
}

func (this *WorkflowManager) AddTaskHandle(t TaskHandle) {
	if t != nil {
		t.SetNotifylistener(this.notifylistener)
		this.taskmanager.addHandle(t)
	}
}

func (this *WorkflowManager) Start(flowname string, parameters ...Parameter) error {
	flow := this.FlowMeta.GetFlow(flowname)
	if flow != nil {
		finstance := Engine{}
		this.instanceNo++
		finstance.Init(fmt.Sprintf("%s_%s_%d", this.Prefix, flow.Name, this.instanceNo), flow)
		this._instance.SetDefault(finstance.Id, &finstance)
		finstance.Start()
		return nil
	}
	return fmt.Errorf("%s flow define not exist!", flowname)
}

func (this *WorkflowManager) notifylistener(flowid string, taskid int, success bool, err string, parameter ...Parameter) {
	fin, exist := this._instance.Get(flowid)
	if exist {
		finstance := fin.(*Engine)
		finstance.CommitTask(taskid, success, err, "", parameter...)

	} else {
		//加载实例

	}

}
