package rpcmock

import (
	"reflect"

	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func NewClient() *Client {
	return &Client{}
}

// Call implements the jsonrpc.Client interface.
func (c *Client) Call(method string, params, result interface{}) error {
	args := c.Called(method, params, result)
	return args.Error(0)
}

func (c *Client) MockResponse(method string, params, response interface{}) {
	if params == nil {
		params = mock.Anything
	}

	call := c.On("Call", method, params, mock.Anything)
	call.Once()
	call.Run(func(args mock.Arguments) {
		reflect.ValueOf(args.Get(2)).Elem().Set(reflect.ValueOf(response))
	})
	call.Return(nil)
}

func (c *Client) MockError(method string, params interface{}, err error) {
	if params == nil {
		params = mock.Anything
	}

	call := c.On("Call", method, params, mock.Anything)
	call.Once()
	call.Return(err)
}
