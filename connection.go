package memgo

import (
	"net"
	"strconv"
)

func NewConnection(client *Client) net.Conn {
	address := client.Host + ":" + strconv.Itoa(client.Port)
	conn, err := net.Dial(client.Transport, address)
	if err != nil {
		panic(err)
	}
	return conn
}
