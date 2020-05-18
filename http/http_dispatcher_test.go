package http

import (
	"fmt"
	"github.com/aosfather/bingo_mvc"
	"log"
	"testing"
)

type MyRequest struct {
	Name string `Field:"name"`
}
type MyHandle struct {
	Test  string `mapper:"name(test);url(/test);method(GET);style(HTML)"`
	Test1 string `mapper:"name(test1);url(/test1);method(GET);style(JSON)"`
}

func (this *MyHandle) GetHandles() bingo_mvc.HandleMap {
	result := bingo_mvc.NewHandleMap()
	r := &MyRequest{}
	result.Add("test", this.DoTest, r)
	result.Add("test1", this.DoTest1, r)
	return result
}

func (this *MyHandle) DoTest(a interface{}) interface{} {
	t := a.(*MyRequest)
	log.Println(t.Name)
	return "hello"
}

func (this *MyHandle) DoTest1(a interface{}) interface{} {
	t := a.(*MyRequest)
	log.Println(t.Name)
	return fmt.Sprintf("hello %s", t.Name)
}
func TestHttpDispatcher_Run(t *testing.T) {
	h := HttpDispatcher{}
	h.Port = 8090
	h.Run()
}

func TestRun(t *testing.T) {
	hd := HttpDispatcher{}
	hd.AddRequestMapperBystruct(&MyHandle{"123", "456"})
	hd.Port = 8090
	hd.Run()

}
