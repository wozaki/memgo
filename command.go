package memgo

import (
	"strconv"
	"strings"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

type Command struct {
	name string
	item Item
}

func (c *Command) buildRequest(shouldCompress bool) ([]byte, error) {
	val, err := c.serialize(shouldCompress)
	if err != nil {
		return nil, err
	}

	byteSize := len(val)
	req := []string{c.name, c.item.Key, strconv.FormatUint(uint64(c.item.Flags),10), strconv.Itoa(c.item.Exptime), strconv.Itoa(byteSize)}

	r1 := append([]byte(strings.Join(req, " ")+Newline), val...)
	r2 := append(r1, []byte(Newline)...)
	return r2, nil
}

func (c *Command) serialize(shouldCompress bool) ([]byte, error) {
	val := []byte(c.item.Value)
	if !shouldCompress {
		return val, nil
	}

	compressed, err := compress(val)
	if err != nil {
		return nil, err
	}

	return compressed, nil
}
