package memgo

import (
	"net"
)

//TODO: Define key as specific type
func NewConnection(client *Client, key string) (net.Conn, error) {
	destination, err := client.Destinations.Get(key)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout("tcp", destination, client.Config.connectTimeout())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
