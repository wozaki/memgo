package memgo

import (
	"strconv"
	"strings"
	"bufio"
	"errors"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

const Newline = "\r\n"

var ErrorNotStored = errors.New("memcached returned NOT_STORED")

func handleErrorResponse(response string) error {
	if strings.HasPrefix(response, "CLIENT_ERROR") {
		return errors.New("memcached returned CLIENT_ERROR: " + response)
	} else {
		panic("returned unexpected value: " + response)
	}
}

type Client struct {
	Destinations Destinations
}

func NewClient(destinations []string) Client {
	return Client{Destinations: NewDestinations(destinations)}
}

var DefaultClient = NewClient([]string{"localhost:11211"})

type Item struct {
	Key     string
	Value   string
	Flags   uint32
	Exptime int // TODO: use time.Duration
}

type Response struct {
	Key     string
	Value   string
	Flags   uint32
	CasId   uint64
}

func Set(item Item) error {
	return DefaultClient.Set(item)
}

func Add(item Item) error {
	return DefaultClient.Add(item)
}

func Get(k string) (response *Response, err error) {
	return DefaultClient.Get(k)
}

func Gets(k string) (response *Response, err error) {
	return DefaultClient.Gets(k)
}

func (c *Client) Set(item Item) error {
	return c.store(Command{name: "set", item: item})
}

func (c *Client) Add(item Item) error {
	return c.store(Command{name: "add", item: item})
}

func (c *Client) Get(k string) (response *Response, err error) {
	return c.retrieve(k, "get")
}

func (c *Client) Gets(k string) (response *Response, err error) {
	return c.retrieve(k, "gets")
}

func (c *Client) store(command Command) error {
	conn := NewConnection(c, command.item.Key)
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
		return handleErrorResponse(s)
	}
}

func (c *Client) retrieve(k string, command string) (response *Response, err error) {
	conn := NewConnection(c, k)
	defer conn.Close()

	req := []string{command, k}
	conn.Write([]byte(strings.Join(req, " ") + Newline))

	// The format is here:
	// VALUE <key> <flags> <bytes> [<cas unique>]\r\n
	// <data block>\r\n
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	heads := strings.Split(scanner.Text(), " ")
	switch heads[0] {
	case "END":
		return nil, nil
	case "VALUE":
		flags, err := strconv.ParseUint(heads[2], 10, 32)
		if err != nil {
			return nil, err
		}
		casId := uint64(0)
		if len(heads) > 4 {
			casId, err = strconv.ParseUint(heads[3], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		scanner.Scan()
		return &Response{Key: k, Value: scanner.Text(), Flags: uint32(flags), CasId: uint64(casId)}, nil
	default:
		return nil, handleErrorResponse(heads[0])
	}
}
