package flow

import (
	"fmt"
	"github.com/aosfather/bingo_utils"
	"github.com/aosfather/bingo_utils/files"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type TTaskStatus byte    //任务状态
type TRetryLogic byte    //重试逻辑
type TTimeoutPolicy byte //超时策略
type TTaskType byte      //任务类型
const (
	TS_Schedule  TTaskStatus = 1  //计划
	TS_Completed TTaskStatus = 2  //完成
	TS_Cancel    TTaskStatus = 4  //取消
	TS_Pause     TTaskStatus = 3  //暂停
	TS_Progress  TTaskStatus = 11 //进行中

	//--------重试逻辑--------//

	RL_fixed   TRetryLogic = 18 //重新安排任务
	RL_backoff TRetryLogic = 19 //重新安排
	//--------超时策略--------//
	TP_retry      TTimeoutPolicy = 21 //重试
	TP_timeout_wf TTimeoutPolicy = 22 //工作流标记为超时并终止
	TP_alert_only TTimeoutPolicy = 23 //仅仅作为警告

	//--------任务类型--------//
	TT_normal       TTaskType = 31 //标准任务
	TT_dynamic      TTaskType = 32 //动态任务
	TT_decide       TTaskType = 33 //判断决策任务
	TT_fork         TTaskType = 34 //并行任务
	TT_dynamic_fork TTaskType = 45 //动态分支并行
	TT_join         TTaskType = 46 //合并任务
	TT_wait         TTaskType = 47 //等待
	TT_http         TTaskType = 48 //http调用的任务
	TT_end          TTaskType = 0  //结束任务节点

)

func (this *TTaskType) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var text string
	unmarshal(&text)
	text = strings.ToLower(text)
	switch text {
	case "normal":
		*this = TT_normal
	case "decide":
		*this = TT_decide
	case "end":
		*this = TT_end
	default:
		*this = TT_normal
	}
	return nil
}

type Flow struct {
	Name             string           //工作流程的名称，唯一
	Description      string           //流程描述
	Version          int              //版本号
	OutputParameters string           `yaml:"outputs"`
	InputParameters  string           `yaml:"inputs"` //输入参数列表，用于记录工作流所需要的输入，可选
	Tasks            map[string]*Task //任务定义
	StartTask        string           `yaml:"start"`
}

func (this *Flow) LoadFromYamlFile(f string) {
	if files.IsFileExist(f) {
		data, err := ioutil.ReadFile(f)
		if err == nil {
			err = yaml.Unmarshal(data, this)
		}
		if err != nil {
			bingo_utils.Debug("load flow file error,", err.Error())
		}
	} else {
		bingo_utils.Debug(f, " the file not exist!")
	}
}

func (this *Flow) AddTask(t *Task) {
	if this.Tasks == nil {
		this.Tasks = make(map[string]*Task)
	}
	this.Tasks[t.TaskName] = t
}

//通过任务名称，获取task的定义
func (this *Flow) GetTask(tname string) *Task {
	for _, t := range this.Tasks {
		if t.TaskName == tname {
			return t
		}
	}

	return nil
}

//任务
type Task struct {
	Name            string    `yaml:"task"` //任务类型
	TaskName        string    `yaml:"code"` //任务名称，用于在流程中使用，必须唯一
	Type            TTaskType //任务类型
	Description     string    //任务描述
	Optional        bool      //是否可忽略，如果是true，当失败后，流程将继续.任务的状态反映为COMPLETED_WITH_ERRORS	默认为 false
	InputParameters Context   `yaml:"inputs"` //任务输入定义
	Expr            string    //表达式
	NextTask        []string  `yaml:"nexts"` //后续任务节点
}

type ParameterDefineType byte

const (
	PT_Const ParameterDefineType = 1 //常量
	PT_Refer ParameterDefineType = 2 //引用
)

type ParameterDefine struct {
	Key   string
	Type  ParameterDefineType
	Value string
}

//任务设置定义
type TaskDefine struct {
	Name                   string            //任务类型，唯一
	RetryCount             int               //重试次数
	RetryLogic             TRetryLogic       //重试机制
	TimeoutSeconds         int64             //超时间 单位毫秒ms
	TimeoutPolicy          TTimeoutPolicy    //超时策略
	ResponseTimeoutSeconds int64             //返回超时时间
	InputParameters        []ParameterDefine //输入参数定义
	OutputKeys             []string          //任务输出字段
}

//流程实例
type FlowInstance struct {
}

type Context map[string]interface{} //上下文
func (this *Context) GetAsString(k string) string {
	v := *this
	value := v[k]
	if value != nil {
		if str, ok := value.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", value)
	}
	return ""
}

//任务实例
type TaskInstance struct {
	FlowInstance    string  //流程的实例Id
	Id              int     //任务的实例Id
	TaskName        string  //任务名称
	InstanceContext Context //任务实例上下文，需要记录的对
	T               *Task
}
