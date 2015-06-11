package client

import (
	"testing"

	"github.com/bmizerany/assert"
)

var c *Client

func init() {
	c = New("http://localhost:4000/rpc")
}

func TestCall(t *testing.T) {
	var res map[string]interface{}
	err := c.Call("Coupons.GetById", map[string]string{"id": "trial"}, &res)
	assert.Equal(t, nil, err)
	assert.Equal(t, "trial", res["id"].(string))
}
