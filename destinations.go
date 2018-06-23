package memgo

import (
	"stathat.com/c/consistent"
)

type Destinations struct {
	Hashing *consistent.Consistent
}

func (d *Destinations) Get(key string) (destination string, err error) {
	return d.Hashing.Get(key)
}

func NewDestinations(destinations []string) Destinations {
	c := consistent.New()
	for _, d := range destinations {
		c.Add(d)
	}
	return Destinations{Hashing: c}
}
