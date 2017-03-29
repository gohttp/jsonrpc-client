package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/abursavich/nett"
	rpcjson "github.com/gohttp/rpc/json"
)

// Client.
type Client interface {
	Call(method string, args interface{}, res interface{}) error
}

// Create new Client.
func NewClient(addr string) Client {
	dialer := &nett.Dialer{
		Resolver:  &nett.CacheResolver{TTL: 5 * time.Minute},
		Timeout:   1 * time.Minute,
		KeepAlive: 1 * time.Minute,
	}
	return &client{
		addr: addr,
		http: &http.Client{
			Transport: &http.Transport{
				Dial:                dialer.Dial,
				MaxIdleConnsPerHost: 512,
			},
			Timeout: 10 * time.Minute,
		},
	}
}

type client struct {
	http *http.Client
	addr string
}

// Call RPC method with args.
func (c *client) Call(method string, args interface{}, res interface{}) error {
	buf, err := rpcjson.EncodeClientRequest(method, args)
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
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		var e struct {
			Error string `json:"error"`
		}

		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return fmt.Errorf("jsonrpc: received non json response, with status %d", resp.StatusCode)
		}

		return fmt.Errorf("jsonrpc: %s", e.Error)
	}

	err = rpcjson.DecodeClientResponse(resp.Body, res)
	if err != nil {
		return err
	}

	return nil
}
