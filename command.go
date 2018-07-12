package memgo

import (
	"strconv"
	"strings"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

type Command struct {
	name string
	key string
	value string
	flags int
	exptime int
}

func (c *Command) buildRequest() []byte {
	byteSize := len(c.value)
	req := []string{c.name, c.key, strconv.Itoa(c.flags), strconv.Itoa(c.exptime), strconv.Itoa(byteSize)}
	return []byte(strings.Join(req, " ") + Newline + c.value + Newline)
}
