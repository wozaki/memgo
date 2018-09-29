package memgo

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

const (
	CompressFlag uint16 = 0x1
)

type Client struct {
	Servers Servers
	Config  Config
}

func NewClient(servers []string, config Config) Client {
	return Client{Servers: NewServers(servers, config.connectTimeout()), Config: config}
}

type Key struct {
	body string
}

func hashMD5(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func newKey(value string, namespase string) Key {
	k := url.QueryEscape(value)
	if len(namespase) > 0 {
		k = namespase + ":" + k
	}
	if len(k) > 250 {
		k = hashMD5(k)
	}
	return Key{body: k}
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
	Flags   Flags
	CasId   uint64
}

func Set(item Item) (res *Response, err error)  {
	return DefaultClient.Set(item)
}

func Get(k string) (response *Response, err error) {
	return DefaultClient.Get(k)
}

func (c *Client) Set(item Item) (res *Response, err error) {
	return c.store("set", item)
}

func (c *Client) Add(item Item) (res *Response, err error) {
	return c.store("add", item)
}

func (c *Client) Get(k string) (response *Response, err error) {
	return c.retrieve("get", k)
}

func (c *Client) Gets(k string) (response *Response, err error) {
	return c.retrieve("gets", k)
}

func (c *Client) store(operation string, item Item) (res *Response, err error) {
	return c.request(NewStorageCommand(operation, item, newKey(item.Key, c.Config.Namespace), c.Config.compressThresholdByte()))
}

func (c *Client) retrieve(operation string, key string) (res *Response, err error) {
	return c.request(NewRetrievalCommand(operation, newKey(key, c.Config.Namespace)))
}

func (c *Client) request(command Command) (res *Response, err error) {
	return c.Servers.request(command)
}
