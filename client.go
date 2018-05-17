package memgo

import (
	"fmt"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

type Client struct {
	Host      string
	Port      int
	Transport string
}

var DefaultClient = &Client{Host: "localhost", Port: 11211, Transport: "tcp"}

type Response struct {
	Status string
}

//TODO: define as specific type
func Set(k string, v string) (resp *Response, err error) {
	return DefaultClient.Set(k, v)
}

func (c *Client) Set(k string, v string) (resp *Response, err error) {
	conn := NewConnection(c)
	defer conn.Close()

	// TODO: Use set command
	conn.Write([]byte("stats\n"))
	conn.Write([]byte("quit\n"))

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		panic(err)
	}
	fmt.Println("reply", string(reply))

	var r = &Response{}
	r.Status = k + ":" + v

	return r, nil
}
