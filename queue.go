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

func (q *Queue) Push(rr dns.RR, ip netip.Addr) bool {
	entry := &QueueEntry{rr: rr, ip: ip}
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.maximum_size == len(q.questions) {
		return false
	}
	q.questions = append(q.questions, *entry)
	q.cond.Signal()
	return true
}

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

func (q *Queue) PeekBlocking() dns.RR {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.questions) == 0 {
		q.cond.Wait()
	}
	rr := q.questions[0]
	return rr.rr
}

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
