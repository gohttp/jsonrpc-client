package jsonrpc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gohttp/rpc/json"
)

type Client interface {
	Call(method string, args interface{}, res interface{}) error
}

func NewClient(addr string) Client {
	return &client{
		addr: addr,
	}
}

type client struct {
	http *http.Client
	addr string
}

// Call invokes the RPC method with the given arguments and stores the result
// in the value pointed by `res`.
func (c *client) Call(method string, args interface{}, res interface{}) error {
	buf, err := json.EncodeClientRequest(method, args)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(buf)

	r, err := http.NewRequest("POST", c.addr, body)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")

	if c.http == nil {
		c.http = http.DefaultClient
	}

	resp, err := c.http.Do(r)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received status code %d with status: %s", resp.StatusCode, resp.Status)
	}

	err = json.DecodeClientResponse(resp.Body, res)
	if err != nil {
		return err
	}

	return nil
}
