package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/abursavich/nett"
	rpcjson "github.com/gohttp/rpc/json"
)

type Client interface {
	Call(method string, args interface{}, res interface{}) error
	CallContext(context.Context, string, interface{}, interface{}) error
}

// NewClient creates a new Client that caches DNS responses for 5 minutes, and
// times out RPC requests after 10 minutes each.
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
	http      *http.Client
	addr      string
	userAgent string
}

// CallContext calls the given RPC method with the given arguments. args is
// serialized to JSON before sending to the remote server. The response is
// decoded into res.
func (c *client) Call(method string, args interface{}, res interface{}) error {
	return c.CallContext(context.Background(), method, args, res)
}

// CallContext calls the given RPC method with the given arguments. args is
// serialized to JSON before sending to the remote server. The response is
// decoded into res.
func (c *client) CallContext(ctx context.Context, method string, args interface{}, res interface{}) error {
	buf, err := rpcjson.EncodeClientRequest(method, args)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(buf)

	r, err := http.NewRequest("POST", c.addr, body)
	if err != nil {
		return err
	}
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Accept-Charset", "utf-8")

	if c.userAgent != "" {
		r.Header.Set("User-Agent", c.userAgent)
	}

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
		var e struct {
			Error string `json:"error"`
		}

		if err = json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return fmt.Errorf("jsonrpc: received non json response, with status %d", resp.StatusCode)
		}

		return fmt.Errorf("jsonrpc: %s", e.Error)
	}

	return rpcjson.DecodeClientResponse(resp.Body, res)
}
