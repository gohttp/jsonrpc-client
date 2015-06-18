package jsonrpc

import (
	"bytes"
	"net/http"

	"github.com/gorilla/rpc/json"
)

// Client.
type Client struct {
	http *http.Client
	addr string
}

// Create new Client.
func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
		http: &http.Client{},
	}
}

// Call RPC method with params.
func (c *Client) Call(method string, params, res interface{}) error {
	buf, err := json.EncodeClientRequest(method, params)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(buf)

	r, err := http.NewRequest("POST", c.addr, body)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(r)

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	return json.DecodeClientResponse(resp.Body, res)
}
