package ns

import (
	"sync"

	"codeberg.org/miekg/dns"
)

type Cache struct {
	data  map[string]*dns.Msg
	order []string
	mutex sync.Mutex

	maximum_size int
}

func newCache(maximum_size int) *Cache {
	return &Cache{data: make(map[string]*dns.Msg), order: make([]string, 0), maximum_size: maximum_size}
}

func (c *Cache) Put(rr dns.RR, msg *dns.Msg) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.order) == c.maximum_size {
		delete(c.data, c.order[0])
		c.order = c.order[1:]
	}
	c.data[rr.String()] = msg
	c.order = append(c.order, rr.String())
}

func (c *Cache) Get(rr dns.RR) (*dns.Msg, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	msg, ok := c.data[rr.String()]
	return msg, ok
}
