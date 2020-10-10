package flow

import "fmt"

//任务接口
type ITask interface {
	GetName() string
}

//管理器:任务定义管理器、流程管理器
type TaskManager struct {
	tasks   map[string]ITask       //任务实现
	defines map[string]*TaskDefine //任务定义

}

func (this *TaskManager) Init() {
	if this.defines == nil {
		this.defines = make(map[string]*TaskDefine)
	}

	if this.tasks == nil {
		this.tasks = make(map[string]ITask)
	}
}

func (this *TaskManager) Register(define *TaskDefine) {
	if define != nil {
		this.defines[define.Name] = define
	}
}

func (this *TaskManager) AddTaskImp(t ITask) {
	if t != nil {
		this.tasks[t.GetName()] = t
	}

}

//获取任务，定义和实现
func (this *TaskManager) GetTask(name string) (*TaskDefine, ITask) {
	return nil, nil
}

//获取任务定义
func (this *TaskManager) GetTaskDefine(name string) *TaskDefine {
	if name != "" && this.defines != nil {
		return this.defines[name]
	}

	return nil
}

//获取任务实现
func (this *TaskManager) GetTaskImp(name string) ITask {
	if name != "" && this.tasks != nil {
		return this.tasks[name]
	}

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

type WorkflowMetaManager interface {
	Publish(flow *Flow) error
	GetFlow(name string) *Flow
	Remove(name string) error
}

type WorkflowManager struct {
	FlowMeta   WorkflowMetaManager
	Prefix     string
	instanceNo int
}

func (this *WorkflowManager) Start(flowname string, parameters ...Parameter) error {
	flow := this.FlowMeta.GetFlow(flowname)
	if flow != nil {
		finstance := Engine{}
		this.instanceNo++
		finstance.Init(fmt.Sprintf("%s_%s_%d", this.Prefix, flow.Name, this.instanceNo), flow)
		finstance.Start()
		return nil
	}
	return fmt.Errorf("%s flow define not exist!", flowname)
}
