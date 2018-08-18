package memgo

import (
	"stathat.com/c/consistent"
)

type Servers struct {
	Hashing *consistent.Consistent
}

func (d *Servers) Get(key string) (server string, err error) {
	return d.Hashing.Get(key)
}

func NewServers(servers []string) Servers {
	c := consistent.New()
	for _, d := range servers {
		c.Add(d)
	}
	return Servers{Hashing: c}
}
