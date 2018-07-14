package memgo

import (
	"strconv"
	"strings"
	"bufio"
	"errors"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

var ErrorNotStored = errors.New("memcached returned NOT_STORED")

type Client struct {
	Destinations Destinations
}

func NewClient(destinations []string) Client {
	return Client{Destinations: NewDestinations(destinations)}
}

var DefaultClient = NewClient([]string{"localhost:11211"})

type Response struct {
	Val      string
	Flags    int
	ByteSize int
}

//TODO: define as specific type
func Set(k string, v string, flags int, exptime int) error {
	return DefaultClient.Set(k, v, flags, exptime)
}

func Add(k string, v string, flags int, exptime int) error {
	return DefaultClient.Add(k, v, flags, exptime)
}

func Get(k string) (resp *Response, err error) {
	return DefaultClient.Get(k)
}

const Newline = "\r\n"

func (c *Client) store(command Command) error {
	conn := NewConnection(c, command.key)
	defer conn.Close()
	conn.Write(command.buildRequest())

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	s := scanner.Text()
	switch s {
	case "STORED":
		return nil
	case "NOT_STORED":
		return ErrorNotStored
	default:
		panic("returned unexpected value: " + s)
	}
}

func (c *Client) Set(k string, v string, flags int, exptime int) error {
	return c.store(Command{name: "set", key: k, value: v, flags: flags, exptime: exptime})
}

func (c *Client) Add(k string, v string, flags int, exptime int) error {
	return c.store(Command{name: "add", key: k, value: v, flags: flags, exptime: exptime})
}

func (c *Client) Get(k string) (resp *Response, err error) {
	conn := NewConnection(c, k)
	defer conn.Close()

	command := "get"

	req := []string{command, k}
	conn.Write([]byte(strings.Join(req, " ") + Newline))

	var r = &Response{}
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	heads := strings.Split(scanner.Text(), " ")
	switch heads[0] {
	case "END":
		return nil, nil
	case "VALUE":
		r.Flags, _ = strconv.Atoi(heads[2])
		r.ByteSize, _ = strconv.Atoi(heads[3])
		scanner.Scan()
		r.Val = scanner.Text()
		return r, nil
	default:
		panic("Unexpected response:" + heads[0])
	}
}
