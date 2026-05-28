package ns

import (
	"fmt"
	"sync"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

func TestCache(t *testing.T) {
	var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}
	cache := newCache(2)

	txt1 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test1"}}}
	txt2 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test2"}}}
	txt3 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test3"}}}

	a1 := []dns.RR{&dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Answer1"}}}}
	a2 := []dns.RR{&dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Answer2"}}}}
	a3 := []dns.RR{&dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Answer3"}}}}

	cache.Put(&txt1, &a1)
	cache.Put(&txt2, &a2)
	cache.Put(&txt3, &a3)

	if len(cache.order) != 2 {
		t.Error("cache length should be 2 as that is the maximum")
		t.Fail()
	}

	if _, b := cache.Get(&txt1); b != false {
		t.Error("cache index txt1 should be nil")
		t.Fail()
	}

	if m, _ := cache.Get(&txt2); (*m)[0].String() != a2[0].String() {
		t.Error("cache index txt2 should be msg2")
		t.Fail()
	}
}

func TestCacheMultiThread(t *testing.T) {

	var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}
	// initialize test data

	cache := newCache(100)

	var txts []dns.RR = make([]dns.RR, 400)
	var anss [][]dns.RR = make([][]dns.RR, 400)
	for i := 0; i < 400; i++ {
		txts[i] = &dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{fmt.Sprintf("Test %d", i)}}}
		anss[i] = []dns.RR{&dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{fmt.Sprintf("Answer %d", i)}}}}
	}

	var wg sync.WaitGroup
	var start sync.WaitGroup
	wg.Add(2)
	start.Add(1)
	// make goroutines that add to the cache
	go func() {
		defer wg.Done()
		start.Wait()
		for i := 0; i < 200; i++ {
			cache.Put(txts[i], &anss[i])
		}
	}()

	go func() {

		defer wg.Done()
		start.Wait()
		for i := 200; i < 400; i++ {
			cache.Put(txts[i], &anss[i])
		}
	}()

	// start goroutines at the same time
	start.Done()

	// wait for them to exit and make some assertions
	wg.Wait()
	if len(cache.order) != cache.maximum_size {
		t.Error("cache order length should be 100 as that is the maximum")
	}

	if len(cache.data) != cache.maximum_size {
		t.Error("cache map length should be 100 as that is the maximum")
	}

	for i := 0; i < 100; i++ {
		if _, b := cache.Get(txts[i]); b != false {
			t.Error(fmt.Sprintf("cache index txt%d should be nil", i))
		}
	}

	for i := 100; i < 200; i++ {
		if m, b := cache.Get(txts[i]); b != false && (*m)[0].String() != anss[i][0].String() {
			t.Error(fmt.Sprintf("cache index txt%d should be nil", i))
		}
	}

	for i := 200; i < 300; i++ {
		if _, b := cache.Get(txts[i]); b != false {
			t.Error(fmt.Sprintf("cache index txt%d should be nil", i))
		}
	}

	for i := 300; i < 400; i++ {
		if m, b := cache.Get(txts[i]); b != false && (*m)[0].String() != anss[i][0].String() {
			t.Error(fmt.Sprintf("cache index txt%d should be nil", i))
		}
	}

}
