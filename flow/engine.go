package flow

import "github.com/aosfather/bingo_utils"

type Parameter struct {
	Key   string
	Value interface{}
}

//流程实例执行引擎
type Engine struct {
	Id              string //实例唯一ID
	Name            string //对应的流程名称
	Version         int    //对应的流程版本
	flow            *Flow
	InstanceContext Context       //流程实例的上下文
	task            *TaskInstance //当前任务节点
	tm              *TaskManager  //任务管理
	instanceNo      int           //实例序号
	taskinstances   map[int]*TaskInstance
}

func (this *Engine) Init(id string, f *Flow) {
	this.flow = f
	this.Id = id
	this.Name = f.Name
	this.Version = f.Version
	this.InstanceContext = make(Context)
	this.taskinstances = make(map[int]*TaskInstance)
}

/*
  开始执行流程
  1、找到启动任务
  2、执行启动任务
*/
func (this *Engine) Start() {
	first := this.flow.GetTask(this.flow.StartTask)
	this.handle(first)

}

func (this *Engine) handle(t *Task) {
	var ti *TaskInstance
	//检查是否还有相同任务在执行，如果在则需要继续处理
	for _, i := range this.taskinstances {
		if i.T == t {
			//如果同一个任务，则继续处理

			ti = i
		}
	}

	//新任务则新建实例
	if ti == nil {
		ti = this.createTaskInstance(t)
	}

	if ti != nil {
		this.taskinstances[ti.Id] = ti
		bingo_utils.Debug(ti)
		switch t.Type {
		case TT_decide:
			this.handle_decide(ti)
		case TT_normal:
			this.handle_normal(ti)
		case TT_end:
			bingo_utils.Debug("finished!")
		}

	}

}

func (this *Engine) handle_decide(ti *TaskInstance) {
	//执行表达式判断
	bingo_utils.Debugf("handle '%v' decide", ti.TaskName)
	//提交当前任务状态，流转到下一个
	this.CommitTask(ti.Id, true, "", ti.T.NextTask[0])

}

func (this *Engine) handle_normal(ti *TaskInstance) {
	bingo_utils.Debugf("handle '%v' normal", ti.TaskName)
	this.CommitTask(ti.Id, true, "", "")
}

func (this *Engine) createTaskInstance(t *Task) *TaskInstance {
	if t != nil {
		this.instanceNo++
		return &TaskInstance{this.Id, this.instanceNo, t.TaskName, make(Context), t}
	}
	return nil
}

//执行任务
func (this *Engine) runTask(name string) {
	_, plugin := this.tm.GetTask(name)
	plugin.GetName()
}

//更新任务状态
func (this *Engine) CommitTask(id int, success bool, msg string, next string, outParameters ...Parameter) {
	if success {
		ti := this.taskinstances[id]
		if ti != nil {
			tname := next
			if tname == "" {
				tname = ti.T.NextTask[0]
			}
			//保存更新任务状态

			delete(this.taskinstances, id)

			t := this.flow.Tasks[tname]
			if t != nil {
				this.handle(t)
			}
		}
	}
}
