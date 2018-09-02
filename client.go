package memgo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

const (
	Newline             = "\r\n"
	CompressFlag uint16 = 0x1
)

type Client struct {
	Servers Servers
	Config  Config
}

func NewClient(servers []string, config Config) Client {
	return Client{Servers: NewServers(servers, config.connectTimeout()), Config: config}
}

var DefaultClient = NewClient([]string{"localhost:11211"}, Config{})

type Flags struct {
	Value   uint16 // use 16 bit for the backward compatibility. In memcached 1.2.1 and higher, flags may be 32-bits.
}

func (f *Flags) shouldCompress() bool {
	return f.Value & CompressFlag != 0
}

type Item struct {
	Key     string
	Value   string
	Flags   Flags
	Exptime int // TODO: use time.Duration
}

type Response struct {
	Key     string
	Value   string
	Flags   uint16
	CasId   uint64
}

func Set(item Item) error {
	return DefaultClient.Set(item)
}

func Get(k string) (response *Response, err error) {
	return DefaultClient.Get(k)
}

func (c *Client) Set(item Item) error {
	return c.store(NewCommand("set", item, c.Config.compressThresholdByte()))
}

func (c *Client) Add(item Item) error {
	return c.store(NewCommand("add", item, c.Config.compressThresholdByte()))
}

func (c *Client) Get(k string) (response *Response, err error) {
	return c.retrieve(k, "get")
}

func (c *Client) Gets(k string) (response *Response, err error) {
	return c.retrieve(k, "gets")
}

func (c *Client) store(command Command) error {
	conn, err := c.Servers.connect(command.item.Key)
	if err != nil {
		return err
	}

	defer conn.Close()
	req, err := command.buildRequest()
	if err != nil {
		return err
	}
	conn.Write(req)

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
	conn, err := c.Servers.connect(k)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	req := []string{command, k}
	conn.Write([]byte(strings.Join(req, " ") + Newline))

	// The format is here:
	// VALUE <key> <flags> <bytes> [<cas unique>]\r\n
	// <data block>\r\n
	bufReader := bufio.NewReader(conn)
	headBytes, _, err := bufReader.ReadLine()
	if err != nil {
		return nil, err
	}

	heads := strings.Split(string(headBytes), " ")
	switch heads[0] {
	case "END":
		return nil, nil
	case "VALUE":
		flags, err := strconv.ParseUint(heads[2], 16, 16)
		if err != nil {
			return nil, err
		}

		byteSize, err := strconv.Atoi(heads[3])
		if err != nil {
			return nil, err
		}
		casId := uint64(0)
		if len(heads) > 4 {
			casId, err = strconv.ParseUint(heads[4], 10, 64)
			if err != nil {
				return nil, err
			}
		}

		// Scanner can't read large data. https://golang.org/pkg/bufio/#Scanner >Scanning stops unrecoverably at EOF, the first I/O error, or a token too large to fit in the buffer
		var buf bytes.Buffer
		written, err := io.CopyN(&buf, bufReader, int64(byteSize))
		if written != int64(byteSize) {
			return nil, fmt.Errorf("cannot read all value: expected %d, actual %d", byteSize, written)
		}
		if err != nil {
			return nil, err
		}

		var val = ""
		if uint16(flags) & CompressFlag != 0 {
			val, err = decompress(buf.Bytes())
			if err != nil {
				return nil, err
			}
		} else {
			val = string(buf.Bytes())
		}
		return &Response{Key: k, Value: val, Flags: uint16(flags), CasId: uint64(casId)}, nil
	default:
		return nil, handleErrorResponse(heads[0])
	}
}
