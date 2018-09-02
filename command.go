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

func (c *Command) buildRequest(config Config) ([]byte, error) {
	val, flags, err := c.serialize(config)
	if err != nil {
		return nil, err
	}

	byteSize := len(val)
	req := []string{c.name, c.item.Key, strconv.FormatUint(uint64(flags.Value),10), strconv.Itoa(c.item.Exptime), strconv.Itoa(byteSize)}

	r1 := append([]byte(strings.Join(req, " ")+Newline), val...)
	r2 := append(r1, []byte(Newline)...)
	return r2, nil
}

func (c *Command) serialize(config Config) ([]byte, Flags, error) {
	val := []byte(c.item.Value)

	if !c.item.Flags.shouldCompress() && len(val) < config.compressThresholdByte() {
		return val, Flags{}, nil
	}

	compressed, err := compress(val)
	if err != nil {
		return nil, Flags{}, err
	}

	return compressed, Flags{Value: CompressFlag}, nil
}
