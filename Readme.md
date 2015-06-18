# jsonrpc-client

Simple rpc client for json-rpc/http.

View the [docs](http://godoc.org/github.com/gohttp/rpc-logger).

# Example

``` go
jsonrpc.New("http://localhost:4000/rpc")
err := c.Call("Coupon.GetById", map[string]string{"id": "trial"}, &res)
```

# License

MIT
