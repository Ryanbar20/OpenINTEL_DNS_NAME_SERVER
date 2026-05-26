package main

import (
	"sync"

	"codeberg.org/miekg/dns"
)

type Cache struct {
	data  map[string]*dns.Msg
	mutex sync.Mutex
}

func newCache() *Cache {
	return &Cache{data: make(map[string]*dns.Msg)}
}

func (c *Cache) put(rr dns.RR, msg *dns.Msg) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[rr.String()] = msg.Copy()
}

func (c *Cache) get(rr dns.RR) (*dns.Msg, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	msg, ok := c.data[rr.String()]
	return msg, ok
}
