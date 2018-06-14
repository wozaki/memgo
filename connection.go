package memgo

import (
	"net"
)

//TODO: Define key as specific type
func NewConnection(client *Client, key string) net.Conn {
	if len(client.Destinations) == 1 {
		address := client.Destinations[0].Address()
		conn, err := net.Dial(client.Transport, address)
		if err != nil {
			panic(err)
		}
		return conn
	} else {
		//TODO support sharding
		panic("not support sharding yet")
	}
}
