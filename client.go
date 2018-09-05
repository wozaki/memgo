package memgo

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
	return c.request(NewStorageCommand("set", item, c.Config.compressThresholdByte()))
}

func (c *Client) Add(item Item) (res *Response, err error) {
	return c.request(NewStorageCommand("add", item, c.Config.compressThresholdByte()))
}

func (c *Client) Get(k string) (response *Response, err error) {
	return c.request(NewRetrievalCommand("get", k))
}

func (c *Client) Gets(k string) (response *Response, err error) {
	return c.request(NewRetrievalCommand("gets", k))
}

func (c *Client) request(command Command) (res *Response, err error) {
	return c.Servers.request(command)
}
