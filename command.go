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

func (c *Command) buildRequest() []byte {
	byteSize := len(c.item.Value)
	req := []string{c.name, c.item.Key, strconv.FormatUint(uint64(c.item.Flags),10), strconv.Itoa(c.item.Exptime), strconv.Itoa(byteSize)}
	return []byte(strings.Join(req, " ") + Newline + c.item.Value + Newline)
}
