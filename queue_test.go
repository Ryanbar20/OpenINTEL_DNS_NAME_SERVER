package ns

import (
	"crypto/rand"
	"fmt"
	"net/netip"
	"sync"
	"testing"
	"time"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

// tests the ordering of the Queue in a multithreaded setting.
func TestMultiThreadQueueOrder(t *testing.T) {
	queue_len := 400
	var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}
	q := newQueue(queue_len)

	var txts []dns.RR = make([]dns.RR, queue_len)
	for i := range queue_len {
		txts[i] = &dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{fmt.Sprintf("Test %d", i)}}}
	}

	var wg sync.WaitGroup
	var start sync.WaitGroup
	wg.Add(2)
	start.Add(1)
	// make goroutines that add to and pop from the queue

	var questionsPopped []dns.RR = make([]dns.RR, 0)
	go func() {

		defer wg.Done()
		start.Wait()
		for i := range queue_len {
			peek := q.PeekBlocking()

			if q.FindRR(peek) != 0 {
				t.Error("Queue Peek should return the first item and not pop the queue")
			}

			questionsPopped = append(questionsPopped, q.PopBlocking())
			if q.FindRR(questionsPopped[i]) >= 0 {
				t.Error("Queue pop should remove item")
			}
		}
	}()

	go func() {
		defer wg.Done()
		start.Wait()
		for i := range queue_len {
			q.Push(txts[i], netip.Addr{})
		}
	}()

	// start goroutines at the same time
	start.Done()
	wg.Wait()

	if len(q.questions) != 0 {
		t.Error("All questions should have been popped already")
	}

	// check that the queue items were received in order
	for i := range queue_len {
		if questionsPopped[i].String() != txts[i].String() {
			t.Error("Queue pops return incorrect order")

		}
	}
}

// Tests the maintaining of the maximum queue size
func TestMaxQueueSize(t *testing.T) {
	var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}
	q := newQueue(5)

	var txts []dns.RR = make([]dns.RR, 10)
	var ips []netip.Addr = make([]netip.Addr, 10)
	for i := range 10 {
		txts[i] = &dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{fmt.Sprintf("Test %d", i)}}}
		ip_bytes := make([]byte, 4)
		rand.Read(ip_bytes)
		ips[i] = netip.AddrFrom4([4]byte(ip_bytes))
	}

	for i, txt := range txts {
		r := q.Push(txt, ips[i])
		if r == -1 && i < 5 {
			t.Error("Queue push should only return -1 if its full")
		}
	}

	for i, ip := range ips[:5] { // only the first 5 entries are in the queue
		qi := q.FindIP(ip)
		if qi != i {
			t.Error("Find IP should return the correct address of the ip in the queue")
		}
	}

	if q.FindIP(ips[9]) != -1 {
		t.Error("Find ip should return -1 on ips that are not in the queue")
	}

}

// tests the Peek and Pop functions for the Queue
func TestPeekAndPopBlocking(t *testing.T) {
	q := newQueue(2)
	var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}
	txt1 := &dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test 1"}}}
	txt2 := &dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test 2"}}}

	var peek dns.RR
	var pop dns.RR
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		peek = q.PeekBlocking()
		wg.Done()

	}()

	go func() {
		pop = q.PopBlocking()
		wg.Done()
	}()
	time.Sleep(1 * time.Second)
	q.Push(txt1, netip.Addr{})
	q.Push(txt2, netip.Addr{})
	wg.Wait()

	if pop != txt1 {
		t.Error("Pop should always remove the record from the queue and Peek should never do that")
	}
	if peek != txt1 && peek != txt2 {
		t.Error(`Peek should always return the first entry from the queue. \n
		This fail means that Peek did not return a queue entry at all`)
	}

}
