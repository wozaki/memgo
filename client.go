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

type Item struct {
	Key     string
	Value   string
	Flags   uint32
	Exptime int // TODO: use time.Duration
}

func Set(item Item) error {
	return DefaultClient.Set(item)
}

func Add(item Item) error {
	return DefaultClient.Add(item)
}

func Get(k string) (item *Item, err error) {
	return DefaultClient.Get(k)
}

const Newline = "\r\n"

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
		if strings.HasPrefix(s, "CLIENT_ERROR") {
			return errors.New("memcached returned CLIENT_ERROR: " + s)
		} else {
			panic("returned unexpected value: " + s)
		}
	}
}

func (c *Client) Set(item Item) error {
	return c.store(Command{name: "set", item: item})
}

func (c *Client) Add(item Item) error {
	return c.store(Command{name: "add", item: item})
}

func (c *Client) Get(k string) (item *Item, err error) {
	conn := NewConnection(c, k)
	defer conn.Close()

	command := "get"

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
		flags, _ := strconv.ParseUint(heads[2], 10, 32)
		scanner.Scan()
		return &Item{Key: k, Value: scanner.Text(), Flags: uint32(flags)}, nil
	default:
		if strings.HasPrefix(heads[0], "CLIENT_ERROR") {
			return nil, errors.New("memcached returned CLIENT_ERROR: " + heads[0])
		} else {
			panic("returned unexpected value: " + heads[0])
		}
	}
}
