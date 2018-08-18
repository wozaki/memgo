package memgo

import (
	"net"
)

//TODO: Define key as specific type
func NewConnection(client *Client, key string) (net.Conn, error) {
	server, err := client.Servers.Get(key)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout("tcp", server, client.Config.connectTimeout())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
