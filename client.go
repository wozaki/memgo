package memgo

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

type Client struct {
	Host string
	Port int
}

var DefaultClient = &Client{Host: "localhost", Port: 11211}

type Response struct {
	Status string
}

//TODO: define as specific type
func Set(k string, v string) (resp *Response, err error) {
	return DefaultClient.Set(k, v)
}

func (c *Client) Set(k string, v string) (resp *Response, err error) {
	var r = &Response{}
	r.Status = k + ":" + v

	return r, nil
}
