package client

import (
	"bytes"
	"net/http"

	"github.com/gorilla/rpc/json"
)

// Client.
type Client struct {
	http *http.Client
}

// Create new Client.
func New() *Client {
	return &Client{
		http: &http.Client{},
	}
}

// Call RPC method with params.
func (c *Client) Call(method string, params, res interface{}) error {
	buf, _ := json.EncodeClientRequest(method, params)
	body := bytes.NewBuffer(buf)
	r, _ := http.NewRequest("POST", "http://localhost:4000/", body)
	r.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(r)

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	return json.DecodeClientResponse(resp.Body, res)
}
