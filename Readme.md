# jsonrpc-client

Simple rpc client for json-rpc/http.

View the [docs](http://godoc.org/github.com/gohttp/jsonrpc-client).

# Example

``` go
c := jsonrpc.NewClient("http://localhost:4000/rpc")
err := c.Call("Coupon.GetById", map[string]string{"id": "trial"}, &res)
```

``` go
c, _ := jsonrpc.NewClientWithOptions("http://localhost:4000/rpc", jsonrpc.UserAgent("rpc_bot/1.0"))
err := c.Call("Coupon.Apply", map[string]string{"account_id": "abcd1234", "coupon_id": "trial"}, &res)
```

# License

MIT
