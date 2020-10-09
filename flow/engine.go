package flow


type HandleTask func(id string,t *Task)
//引擎
type Engine struct {
	instance FlowInstance //实例
	flow     *Flow
	tm       *TaskManager //任务管理
}

func (this *Engine) Init(id string, f *Flow) {
	this.flow = f
	this.instance = FlowInstance{Id:id,Version: f.Version,Name: f.Name}
	t:=this.flow.Tasks[0]
	taskInstance:=TaskInstance{id,"",t.Name}
	this.instance.task=&taskInstance

}

//执行
func (this *Engine) Run() {
RUNTASK:
	//获取当前节点
	taskinstance := this.instance.task
	//获取节点定义，产生任务，执行
	t := this.flow.GetTask(taskinstance.TaskName)
	this.runTask(t.Name)
	//更新任务状态
	this.commitTask()
	//更新流程节点
	this.instance.task = this.GetNext()
	//如果没有节点，表示已经进行完成
	if this.instance.task.Id == "" {
		//TODO 更新流程状态为完成
		return
	} else {
		//执行下一步
		goto RUNTASK

	}

}

//执行任务
func (this *Engine) runTask(name string) {
	_, plugin := this.tm.GetTask(name)
	plugin.GetName()
}

//更新任务状态
func (this *Engine) commitTask() {

}

//获取下一个节点
func (this *Engine) GetNext() TaskInstance {
	return TaskInstance{}
}
