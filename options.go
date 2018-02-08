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

// HTTPClient sets the http.Client on the Client for the connection.
// If RoundTripper was previously specified, this function will override
// the transport.
func HTTPClient(hc *http.Client) optionFunc {
	return optionFunc(func(c *client) error {
		c.http = hc
		return nil
	})
}

// RoundTripper sets the Transport on the http.Client if the http client has
// already been specified, this overrides the transport.  If not, we create a
// new client.
func RoundTripper(rt http.RoundTripper) optionFunc {
	return optionFunc(func(c *client) error {
		if c.http == nil {
			c.http = &http.Client{
				Transport: rt,
				Timeout:   10 * time.Minute,
			}
		} else {
			c.http.Transport = rt
		}
		return nil
	})
}

// UserAgent sets the User-Agent header on outgoing requests.
func UserAgent(userAgent string) optionFunc {
	return optionFunc(func(c *client) error {
		c.userAgent = userAgent
		return nil
	})
}
