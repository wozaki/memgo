package memgo

import (
	"net"
	"stathat.com/c/consistent"
	"time"
)

type Servers struct {
	Hashing        *consistent.Consistent
	ConnectTimeout time.Duration
}

func NewServers(servers []string, connectTimeout time.Duration) Servers {
	c := consistent.New()
	for _, d := range servers {
		c.Add(d)
	}
	return Servers{Hashing: c, ConnectTimeout: connectTimeout}
}

func (s *Servers) request(command Command) (res *Response, err error) {
	conn, err := s.connect(command.Key().body)
	if err != nil {
		return nil, err
	}

	res, err = command.Perform(conn)
	closedErr := conn.Close()
	if closedErr != nil {
		return nil, closedErr
	}

	return res, err
}

func (s *Servers) connect(key string) (net.Conn, error) {
	server, err := s.Hashing.Get(key)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout("tcp", server, s.ConnectTimeout)
	if err != nil {
		return nil, err
	}

	return conn, nil
}