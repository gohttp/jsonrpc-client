package jsonrpc

import (
	"context"
	"fmt"
)

func Example() {
	c := NewClient("http://localhost:5000/rpc")
	type User struct {
		Email string
	}
	u := new(User)
	err := c.CallContext(context.TODO(), "User.GetOne", "test@example.com", u)
	fmt.Println(err)
	fmt.Println(u.Email)
}
