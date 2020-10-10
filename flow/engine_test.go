package flow

import (
	"github.com/aosfather/bingo_utils"
	"testing"
	"time"
)

func TestEngine_Start(t *testing.T) {
	bingo_utils.SetLogDebugFunc(func(msg ...interface{}) { t.Log(msg...) })
	t1 := &Task{TaskName: "first", Type: TT_normal, NextTask: []string{"second"}}
	t2 := &Task{TaskName: "second", Type: TT_decide, NextTask: []string{"thirt", "first"}}
	t3 := &Task{TaskName: "thirt", Type: TT_normal, NextTask: []string{"first"}}
	f1 := &Flow{Name: "test"}
	f1.AddTask(t1)
	f1.AddTask(t2)
	f1.AddTask(t3)
	f1.StartTask = "first"
	engine := Engine{}
	engine.Init("111", f1)
	engine.Start()
}

func TestWorkflowManager_Start(t *testing.T) {
	bingo_utils.SetLogDebugFunc(func(msg ...interface{}) { t.Log(msg...) })
	f := Flow{}
	f.LoadFromYamlFile("../exampleflow.yaml")
	bingo_utils.Debug(f)
	meta := FlowManager{}
	meta.Init()
	meta.Publish(&f)
	wfm := WorkflowManager{FlowMeta: &meta, Prefix: "wf"}
	wfm.Start("example")
	n := time.Now()
	for i := 1; i < 1000; i++ {
		wfm.Start("example")
	}

	e := time.Now()
	t.Log(e.Sub(n).String())
}
