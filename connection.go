package memgo

import (
	"net"
)

//TODO: Define key as specific type
func NewConnection(client *Client, key string) net.Conn {
	destination, err := client.Destinations.Get(key)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTimeout("tcp", destination, client.connectTimeout())
	if err != nil {
		panic(err)
	}

	return conn
}
