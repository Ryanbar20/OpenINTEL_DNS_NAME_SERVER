package ns

import (
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

const dom = "test.nl."

var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}

func TestCache(t *testing.T) {
	cache := newCache(2)

	txt1 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test1"}}}
	txt2 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test2"}}}
	txt3 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test3"}}}

	a1 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Answer1"}}}
	a2 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Answer2"}}}
	a3 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Answer3"}}}

	msg1 := dns.NewMsg(txt1.String(), dns.TypeTXT)
	msg1.Answer = append(msg1.Answer, &a1)
	msg2 := dns.NewMsg(txt2.String(), dns.TypeTXT)
	msg2.Answer = append(msg2.Answer, &a2)
	msg3 := dns.NewMsg(txt3.String(), dns.TypeTXT)
	msg3.Answer = append(msg3.Answer, &a3)

	cache.put(&txt1, msg1)
	cache.put(&txt2, msg2)
	cache.put(&txt3, msg3)

	if len(cache.order) != 2 {
		t.Error("cache length should be 2 as that is the maximum")
		t.Fail()
	}

	if _, b := cache.get(&txt1); b != false {
		t.Error("cache index txt2 should be msg2")
		t.Fail()
	}

	if m, _ := cache.get(&txt2); m.String() != msg2.String() {
		t.Error("cache index txt2 should be msg2")
		t.Fail()
	}
}
