package ns

import (
	"fmt"
	"net/netip"
	"sync"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

func TestQueueOrder(t *testing.T) {
	var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}
	q := newQueue(400)

	var txts []dns.RR = make([]dns.RR, 400)
	for i := range 400 {
		txts[i] = &dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{fmt.Sprintf("Test %d", i)}}}
	}

	var wg sync.WaitGroup
	var start sync.WaitGroup
	wg.Add(2)
	start.Add(1)
	// make goroutines that add to the queue
	go func() {
		defer wg.Done()
		start.Wait()
		for i := range 400 {
			q.Push(txts[i], netip.Addr{})
		}
	}()

	var questionsPopped []dns.RR = make([]dns.RR, 0)
	go func() {

		defer wg.Done()
		start.Wait()
		for i := range 400 {
			questionsPopped = append(questionsPopped, q.PopBlocking())
			if q.FindRR(questionsPopped[i]) >= 0 {
				t.Error("Queue pop should remove item")
			}
		}
	}()

	// start goroutines at the same time
	start.Done()
	wg.Wait()

	if len(q.questions) != 0 {
		t.Error("All questions should have been popped already")
	}

	// check that the queue items were received in order
	for i := range 400 {
		if questionsPopped[i].String() != txts[i].String() {
			t.Error("Queue pops return incorrect order")

		}
	}
}
