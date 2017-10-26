package rpcmock

import (
	"context"
	"errors"
	"testing"

	"github.com/bmizerany/assert"
)

func TestMockResponse(t *testing.T) {
	c := NewClient()

	c.MockResponse("RPC.Test", "hello", "world")
	c.MockResponse("RPC.Test2", "bye", "world")

{
	var response string
	err := c.Call("RPC.Test", "hello", &response)
	assert.Equal(t, "world", response)
	assert.Equal(t, nil, err)
}
{
	var response string
	err := c.Call("RPC.Test2", "bye", &response)
	assert.Equal(t, "world", response)
	assert.Equal(t, nil, err)
}
	c.AssertExpectations(t)
}

func TestMock(t *testing.T) {
	c := NewClient()

	c.MockError("RPC.Test", "hello", errors.New("somebody set up the bomb"))

	var response string
	err := c.Call("RPC.Test", "hello", &response)

	assert.Equal(t, "", response)
	assert.Equal(t, errors.New("somebody set up the bomb"), err)
	c.AssertExpectations(t)
}

func TestCallContextError(t *testing.T) {
	c := NewClient()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := c.CallContext(ctx, "RPC.Test", "hello", nil)
	assert.Equal(t, err, context.Canceled)
}
