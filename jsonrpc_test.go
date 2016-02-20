package jsonrpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type Args struct {
	A, B int
}

type Arith int

type Result int

func (t *Arith) Multiply(r *http.Request, args *Args, result *Result) error {
	*result = Result(args.A * args.B)
	return nil
}

func TestMockResponse(t *testing.T) {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	arith := new(Arith)
	s.RegisterService(arith, "")

	ts := httptest.NewServer(s)
	defer ts.Close()

	client := NewClient(ts.URL)
	var result int
	err := client.Call("Arith.Multiply", Args{2, 3}, &result)

	assert.Equal(t, nil, err)
	assert.Equal(t, 6, result)
}
