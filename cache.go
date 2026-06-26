package ns

import (
	"sync"

	"codeberg.org/miekg/dns"
)

// A Cache struct for storing database query results
type Cache struct {
	data  map[string][]dns.RR
	order []*string
	mutex sync.Mutex

	maximum_size int
}

// Create a new cache with a maximum size
func newCache(maximum_size int) *Cache {
	return &Cache{data: make(map[string][]dns.RR), order: make([]*string, 0), maximum_size: maximum_size}
}

// puts an entry into the cache. Automatically maintains maximum cache size
func (c *Cache) Put(rr dns.RR, rrs *[]dns.RR) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.order) == c.maximum_size {
		delete(c.data, *c.order[0])
		c.order = c.order[1:]
	}
	str := rr.String()
	c.data[str] = *rrs
	c.order = append(c.order, &str)
}

// gets an entry from the cache
func (c *Cache) Get(rr dns.RR) (*[]dns.RR, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	rrs, ok := c.data[rr.String()]
	return &rrs, ok
}
