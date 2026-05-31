package ns

import (
	"net/netip"
	"sync"

	"codeberg.org/miekg/dns"
)

type QueueEntry struct {
	rr dns.RR
	ip netip.Addr
}

type Queue struct {
	questions []QueueEntry
	cond      sync.Cond

	maximum_size int
}

func newQueue(maximum_size int) *Queue {

	return &Queue{questions: make([]QueueEntry, 0), cond: *sync.NewCond(&sync.Mutex{}), maximum_size: maximum_size}
}

// pushes an entry into the cache
func (q *Queue) Push(rr dns.RR, ip netip.Addr) int {
	entry := &QueueEntry{rr: rr, ip: ip}
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	defer q.cond.Signal()
	if q.maximum_size == len(q.questions) {
		return -1
	}
	q.questions = append(q.questions, *entry)
	return len(q.questions) - 1 // return the index of the question
}

// finds the index of the question RR. -1 if not found
func (q *Queue) FindRR(rr dns.RR) int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for i, question := range q.questions {
		if question.rr.String() == rr.String() {
			return i
		}
	}
	return -1
}

// finds the index of the IP. -1 if not found
func (q *Queue) FindIP(ip netip.Addr) int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for i, question := range q.questions {
		if question.ip == question.ip {
			return i
		}
	}
	return -1
}

// waits until it can read the first entry of the queue
func (q *Queue) PeekBlocking() dns.RR {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.questions) == 0 {
		q.cond.Wait()
	}
	rr := q.questions[0]
	return rr.rr
}

// waits until it can pop the first value off the queue
func (q *Queue) PopBlocking() dns.RR {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.questions) == 0 {
		q.cond.Wait()
	}
	rr := q.questions[0]
	q.questions = q.questions[1:]
	return rr.rr
}
