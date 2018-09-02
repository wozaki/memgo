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

func (s *Servers) Get(key string) (server string, err error) {
	return s.Hashing.Get(key)
}

func NewServers(servers []string, connectTimeout time.Duration) Servers {
	c := consistent.New()
	for _, d := range servers {
		c.Add(d)
	}
	return Servers{Hashing: c, ConnectTimeout: connectTimeout}
}

func (s *Servers) connect(key string) (net.Conn, error) {
	server, err := s.Get(key)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout("tcp", server, s.ConnectTimeout)
	if err != nil {
		return nil, err
	}

	return conn, nil
}