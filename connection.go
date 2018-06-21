package memgo

import (
	"net"
)

//TODO: Define key as specific type
func NewConnection(client *Client, key string) net.Conn {
	address, err := client.Destinations.GetAddress(key)
	if err != nil {
		panic(err)
	}

	conn, err := net.Dial(client.Transport, address)
	if err != nil {
		panic(err)
	}

	return conn
}
