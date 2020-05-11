package http

import (
	"github.com/aosfather/bingo_mvc"
	"testing"
)

type MyRequest struct {
	Name string `Field:"name"`
}
type MyHandle struct {
	Test  string `mapper:"name(test);url(/test);method(GET);style(JSON)"`
	Test1 string `mapper:"name(test1);url(/test1);method(GET);style(JSON)"`
}

func (this *MyHandle) GetHandles() map[string]bingo_mvc.CMap {
	result := make(map[string]bingo_mvc.CMap)
	result["test"] = bingo_mvc.CMap{this.DoTest, &MyRequest{}}
	result["test1"] = bingo_mvc.CMap{this.DoTest1, &MyRequest{}}
	return result
}

func (this *MyHandle) DoTest(a interface{}) interface{} {
	return "hello"
}

func (this *MyHandle) DoTest1(a interface{}) interface{} {
	return "hello1"
}
func TestHttpDispatcher_Run(t *testing.T) {
	h := HttpDispatcher{}
	h.Port = 8090
	h.Run()
}

func TestRun(t *testing.T) {
	hd := HttpDispatcher{}
	hd.AddRequestMapperBystruct(&MyHandle{"123", "456"}, &MyRequest{}, &MyRequest{})
	hd.Port = 8090
	hd.Run()

}
