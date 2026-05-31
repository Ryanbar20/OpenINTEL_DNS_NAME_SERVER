package ns

import (
	"net"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

// mock for a dns Response writer
type mockDnsResponseWriter struct {
	dns.ResponseWriter
}

func (m *mockDnsResponseWriter) WriteMsg(msg *dns.Msg) error {
	return nil
}

func (m *mockDnsResponseWriter) RemoteAddr() net.Addr {
	return &net.IPAddr{}
}

func (m *mockDnsResponseWriter) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (m *mockDnsResponseWriter) Conn() net.Conn {
	return nil
}

func TestShouldRefuse(t *testing.T) {

	// invalid Name format
	REFtxtRR1 := dns.TXT{Hdr: dns.Header{Name: "test.history.openintel.nl.", Class: dns.ClassINET}}
	REFtxtRR2 := dns.TXT{Hdr: dns.Header{Name: "2006c1212.google.com.history.openintel.nl.", Class: dns.ClassINET}}
	REFtxtRR3 := dns.TXT{Hdr: dns.Header{Name: "20061212.google.com.hirtory.openintel.nl.", Class: dns.ClassINET}}
	// valid format and type
	ACCtxtRR := dns.TXT{
		Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET},
	}

	// invalid type
	REFsoaRR := dns.SOA{Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET}}

	REFmsg1 := new(dns.Msg)
	REFmsg1.Question = append(REFmsg1.Question, &REFtxtRR1)
	REFmsg1.Pack()
	REFmsg2 := new(dns.Msg)
	REFmsg2.Question = append(REFmsg2.Question, &REFtxtRR2)
	REFmsg2.Pack()
	REFmsg3 := new(dns.Msg)
	REFmsg3.Question = append(REFmsg3.Question, &REFtxtRR3)
	REFmsg3.Pack()
	REFmsg4 := new(dns.Msg)
	REFmsg4.Pack()

	REFmsg5 := new(dns.Msg)
	REFmsg5.Question = append(REFmsg5.Question, &REFsoaRR)
	REFmsg5.Pack()
	ACCmsg6 := new(dns.Msg)
	ACCmsg6.Question = append(ACCmsg6.Question, &ACCtxtRR)
	ACCmsg6.Pack()

	ns := NewNameServer(10, 10)
	w := &mockDnsResponseWriter{}

	if ns.handle(nil, w, REFmsg1); REFmsg1.Rcode != dns.RcodeRefused {
		t.Error("message 1 should be refused as it has invalid name format")

	}
	if ns.handle(nil, w, REFmsg2); REFmsg2.Rcode != dns.RcodeRefused {
		t.Error("message 2 should be refused as it has invalid name format")

	}
	if ns.handle(nil, w, REFmsg3); REFmsg3.Rcode != dns.RcodeRefused {
		t.Error("message 3 should be refused as it has invalid name format")

	}
	if ns.handle(nil, w, REFmsg4); REFmsg4.Rcode != dns.RcodeRefused {
		t.Error("message 4 should be refused as it has no questions")

	}
	if ns.handle(nil, w, REFmsg5); REFmsg5.Rcode != dns.RcodeRefused {
		t.Error("message 5 should be refused as it has invalid type (soa)")

	}
	if ns.handle(nil, w, ACCmsg6); ACCmsg6.Rcode != dns.RcodeSuccess {
		t.Error("message 6 should be accepted")
	}

}

func TestHandleCacheHit(t *testing.T) {

	txtQRR := dns.TXT{
		Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET},
	}
	txtARR := []dns.RR{&dns.TXT{
		Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET},
		TXT: rdata.TXT{Txt: []string{"Test 2"}},
	}}

	ns := NewNameServer(10, 10)
	ns.cache.Put(&txtQRR, &txtARR)

	msg := new(dns.Msg)
	msg.Question = append(msg.Question, &txtQRR)
	msg.Pack()
	w := &mockDnsResponseWriter{}

	ns.handle(nil, w, msg)
	if msg.Answer[0].String() != txtARR[0].String() {
		t.Error("the name-server should have seen the record in the cache")
	}
}
