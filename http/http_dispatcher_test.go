package http

import "testing"

func TestHttpDispatcher_Run(t *testing.T) {
	h := HttpDispatcher{}
	h.Port = 8090
	h.Run()
}
