package mock

import (
	"errors"
	"testing"

	"github.com/bmizerany/assert"
)

func TestMockResponse(t *testing.T) {
	c := NewClient()

	c.MockResponse("RPC.Test", "hello", "world")

	var response string
	err := c.Call("RPC.Test", "hello", &response)

	assert.Equal(t, "world", response)
	assert.Equal(t, nil, err)
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
