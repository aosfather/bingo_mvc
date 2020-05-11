package fasthttp

import "testing"

func TestFastHTTPDispatcher_Run(t *testing.T) {
	h := FastHTTPDispatcher{}
	h.Port = 8090
	h.Run()

}
