package memgo

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

//TODO: initialize host and port
type Client struct {
}

var DefaultClient = &Client{}

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
