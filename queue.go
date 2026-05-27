package ns

import (
	"sync"

	"codeberg.org/miekg/dns"
)

type Queue struct {
	questions []dns.RR
	cond      sync.Cond

	maximum_size int
}

func newQueue(maximum_size int) *Queue {

	return &Queue{questions: make([]dns.RR, 0), cond: *sync.NewCond(&sync.Mutex{}), maximum_size: maximum_size}
}

func (q *Queue) Push(rr dns.RR) bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.maximum_size == len(q.questions) {
		return false
	}
	q.questions = append(q.questions, rr)
	q.cond.Signal()
	return true
}

func (q *Queue) Find(rr dns.RR) int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for i, question := range q.questions {
		if question.String() == rr.String() {
			return i
		}
	}
	return -1
}

func (q *Queue) PopBlocking() dns.RR {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.questions) == 0 {
		q.cond.Wait()
	}
	rr := q.questions[0]
	q.questions = q.questions[1:]
	return rr
}
