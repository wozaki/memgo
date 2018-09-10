package memgo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

const (
	Newline = "\r\n"
)

type Command interface {
	Perform(conn net.Conn) (res *Response, err error)
	Key() Key
}

type StorageCommand struct {
	name                  string
	key                   Key
	value                 string
	flags                 Flags
	exptime               int
	compressThresholdByte int
}

func NewStorageCommand(name string, item Item, compressThresholdByte int) Command {
	return &StorageCommand{name: name, key: newKey(item.Key), value: item.Value, flags: item.Flags, exptime: item.Exptime, compressThresholdByte: compressThresholdByte}
}

func (c *StorageCommand) Perform(conn net.Conn) (res *Response, err error) {
	req, flags, err := c.buildRequest()
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	s := scanner.Text()
	switch s {
	case "STORED":
		return &Response{Key: c.key.body, Value: c.value, Flags: flags}, nil
	case "NOT_STORED":
		return nil, ErrorNotStored
	default:
		return nil, handleErrorResponse(s)
	}
}

func (c *StorageCommand) Key() Key {
	return c.key
}

func (c *StorageCommand) buildRequest() ([]byte, Flags, error) {
	val, flags, err := c.serialize()
	if err != nil {
		return nil, Flags{}, err
	}

	byteSize := len(val)
	req := []string{c.name, c.key.body, strconv.FormatUint(uint64(flags.Value),10), strconv.Itoa(c.exptime), strconv.Itoa(byteSize)}

	r1 := append([]byte(strings.Join(req, " ")+Newline), val...)
	r2 := append(r1, []byte(Newline)...)
	return r2, flags, nil
}

func (c *StorageCommand) serialize() ([]byte, Flags, error) {
	val := []byte(c.value)

	if !c.shouldCompress(val) {
		return val, Flags{}, nil
	}

	compressed, err := compress(val)
	if err != nil {
		return nil, Flags{}, err
	}

	return compressed, Flags{Value: CompressFlag}, nil
}

func (c *StorageCommand) shouldCompress(value []byte) bool {
	if c.flags.shouldCompress() {
		return true
	}

	if len(value) >= c.compressThresholdByte {
		return true
	}

	return false
}

type RetrievalCommand struct {
	name string
	key Key
}

func NewRetrievalCommand(name string, key string) Command {
	return &RetrievalCommand{name: name, key: newKey(key)}
}

func (c *RetrievalCommand) Perform(conn net.Conn) (res *Response, err error) {
	conn.Write(c.buildRequest())

	// The format is here:
	// VALUE <Key> <flags> <bytes> [<cas unique>]\r\n
	// <data block>\r\n
	bufReader := bufio.NewReader(conn)
	headBytes, _, err := bufReader.ReadLine()
	if err != nil {
		return nil, err
	}

	heads := strings.Split(string(headBytes), " ")
	switch heads[0] {
	case "END":
		return &Response{}, nil
	case "VALUE":
		rawFlags, err := strconv.ParseUint(heads[2], 16, 16)
		if err != nil {
			return nil, err
		}

		flags := Flags{Value: uint16(rawFlags)}
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
		if flags.shouldCompress() {
			val, err = decompress(buf.Bytes())
			if err != nil {
				return nil, err
			}
		} else {
			val = string(buf.Bytes())
		}
		return &Response{Key: c.key.body, Value: val, Flags: flags, CasId: uint64(casId)}, nil
	default:
		return nil, handleErrorResponse(heads[0])
	}
}

func (c *RetrievalCommand) Key() Key {
	return c.key
}

func (c *RetrievalCommand) buildRequest() []byte {
	req := []string{c.name, c.key.body}
	return []byte(strings.Join(req, " ") + Newline)
}
