package memgo

import (
	"strconv"
	"strings"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

type Command struct {
	name string
	item Item
	compressThresholdByte int
}

func NewCommand(name string, item Item, compressThresholdByte int) Command {
	return Command{name: name, item: item, compressThresholdByte: compressThresholdByte}
}

func (c *Command) buildRequest() ([]byte, error) {
	val, flags, err := c.serialize()
	if err != nil {
		return nil, err
	}

	byteSize := len(val)
	req := []string{c.name, c.item.Key, strconv.FormatUint(uint64(flags.Value),10), strconv.Itoa(c.item.Exptime), strconv.Itoa(byteSize)}

	r1 := append([]byte(strings.Join(req, " ")+Newline), val...)
	r2 := append(r1, []byte(Newline)...)
	return r2, nil
}

func (c *Command) serialize() ([]byte, Flags, error) {
	val := []byte(c.item.Value)

	if !c.shouldCompress(val) {
		return val, Flags{}, nil
	}

	compressed, err := compress(val)
	if err != nil {
		return nil, Flags{}, err
	}

	return compressed, Flags{Value: CompressFlag}, nil
}

func (c *Command) shouldCompress(value []byte) bool {
	if c.item.Flags.shouldCompress() {
		return true
	}

	if len(value) >= c.compressThresholdByte {
		return true
	}

	return false
}
