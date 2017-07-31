package jsonrpc

import (
	"net/http"
	"time"
)

type Option interface {
	apply(*client) error
}

func NewClientWithOptions(addr string, options ...Option) (Client, error) {
	c := &client{
		addr: addr,
	}
	for _, o := range options {
		err := o.apply(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

type optionFunc func(*client) error

func (o optionFunc) apply(c *client) error {
	return o(c)
}

func RoundTripper(rt http.RoundTripper) optionFunc {
	return optionFunc(func(c *client) error {
		c.http = &http.Client{
			Transport: rt,
			Timeout:   10 * time.Minute,
		}
		return nil
	})
}

func UserAgent(userAgent string) optionFunc {
	return optionFunc(func(c *client) error {
		c.userAgent = userAgent
		return nil
	})
}
