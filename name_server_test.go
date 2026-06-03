package ns

import (
	"crypto/rand"
	"net"
	"net/netip"
	"testing"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

// mock for a dns Response writer
type mockDnsResponseWriter struct {
	dns.ResponseWriter
	remoteAddr net.Addr
}

func (m *mockDnsResponseWriter) WriteMsg(msg *dns.Msg) error {
	return nil
}

func (m *mockDnsResponseWriter) RemoteAddr() net.Addr {
	if &m.remoteAddr == nil {
		return &net.UDPAddr{}
	}
	return m.remoteAddr
}

func (m *mockDnsResponseWriter) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (m *mockDnsResponseWriter) Conn() net.Conn {
	return nil
}

func (m *mockDnsResponseWriter) setRemoteAddr(addr net.Addr) {
	m.remoteAddr = addr
}

// Tests that the name-server refuses invalid queries or queries for wrong types
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

	ns := NewNameServer(10, 10, "[::]", NAME_SERVER_PORT, MEMORY_LIMIT)
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

// Tests that the name server correctly detects and handles cache-hits.
func TestHandleCacheHit(t *testing.T) {

	txtQRR := dns.TXT{
		Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET},
	}
	txtARR := []dns.RR{&dns.TXT{
		Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET},
		TXT: rdata.TXT{Txt: []string{"Test 2"}},
	}}

	ns := NewNameServer(10, 10, "[::]", NAME_SERVER_PORT, MEMORY_LIMIT)
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

// Tests that the name server correctly handles situations where a user has a query in the queue already (by sending a LIMIT)
func TestHandleIPinQueue(t *testing.T) {

	// initalise two distince records
	txtQRR1 := dns.TXT{
		Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET},
	}
	txtQRR2 := dns.TXT{
		Hdr: dns.Header{Name: "20061222.google.com.history.openintel.nl.", Class: dns.ClassINET},
	}

	// create a random UDP address
	ip_bytes := make([]byte, 4)
	rand.Read(ip_bytes)
	addr := net.UDPAddr{IP: net.IP(ip_bytes)}

	// create a name server and push qrr1 into the queue
	// the queue entry will have the same address as addr
	ns := NewNameServer(10, 10, "[::]", NAME_SERVER_PORT, MEMORY_LIMIT)
	ns.query_queue.Push(&txtQRR1, netip.AddrFrom4([4]byte(ip_bytes)))

	// create a dns message and submit it to the name server under addr
	// the ip is in the queue, but since the question is in the queue aswell, a WAIT message is returned
	{
		msg := new(dns.Msg)
		msg.Question = append(msg.Question, &txtQRR1)
		msg.Pack()
		w := &mockDnsResponseWriter{}
		w.setRemoteAddr(&addr)
		ns.handle(nil, w, msg)
		switch ans := msg.Answer[0].(type) {
		case *dns.HINFO:
			if ans.Cpu != "WAIT" || ans.Os != "You are in queue position 0" {
				t.Error("Query is in the queue, so a wait message should be returned")
			}
		default:
			t.Error("Query is in the queue, so a wait message should be returned")
		}
	}

	// create and submit a new dns message with a different question than the queued one but with the same address (addr)
	// since the user is in the queue and the question is not, a LIMIT message is sent
	{
		msg := new(dns.Msg)
		msg.Question = append(msg.Question, &txtQRR2)
		msg.Pack()
		w := &mockDnsResponseWriter{}
		w.setRemoteAddr(&addr)
		ns.handle(nil, w, msg)
		switch ans := msg.Answer[0].(type) {
		case *dns.HINFO:
			if ans.Cpu != "LIMIT" || ans.Os != "You already have a query in queue position 0" {
				t.Error("User is in the queue and answer is not in cache, so a limit message should be returned")
			}
		default:
			t.Error("User is in the queue and answer is not in cache, so a limit message should be returned")
		}
	}
}

// Tests that the name server correctly handles situations where the query queue is full (by sending a LIMIT)
func TestHandleQueueLimit(t *testing.T) {

	txtQRR1 := dns.TXT{
		Hdr: dns.Header{Name: "20061212.google.com.history.openintel.nl.", Class: dns.ClassINET},
	}
	txtQRR2 := dns.TXT{
		Hdr: dns.Header{Name: "20071212.google.com.history.openintel.nl.", Class: dns.ClassINET},
	}

	ns := NewNameServer(10, 1, "[::]", NAME_SERVER_PORT, MEMORY_LIMIT)
	ip_bytes := make([]byte, 4)
	// create and submit a query from a random IP that gets put into the queue
	{
		msg := new(dns.Msg)
		msg.Question = append(msg.Question, &txtQRR1)
		msg.Pack()
		w := &mockDnsResponseWriter{}
		rand.Read(ip_bytes)
		addr := net.TCPAddr{IP: net.IP(ip_bytes)}
		w.setRemoteAddr(&addr)
		ns.handle(nil, w, msg)
		switch ans := msg.Answer[0].(type) {
		case *dns.HINFO:
			if ans.Cpu != "WAIT" || ans.Os != "You are in queue position 0" {
				t.Error("Query is put in the queue, so a WAIT should be sent as the queue is not over its limit")
			}
		default:
			t.Error("Query is put in the queue, so a WAIT should be sent as the queue is not over its limit")
		}
	}
	// create and submit a query from a different IP that gets put into the queue
	// this adding to the queue sees that the queue limit was reached, and thus a LIMIT message is returned instead
	{
		msg := new(dns.Msg)
		msg.Question = append(msg.Question, &txtQRR2)
		msg.Pack()
		w := &mockDnsResponseWriter{}
		ip_bytes[len(ip_bytes)-1] ^= 0xFF
		addr := net.TCPAddr{IP: net.IP(ip_bytes)}
		w.setRemoteAddr(&addr)
		ns.handle(nil, w, msg)
		switch ans := msg.Answer[0].(type) {
		case *dns.HINFO:
			if ans.Cpu != "LIMIT" || ans.Os != "Queue limit reached" {
				t.Error("Query is put in the queue, but the queue is over its limit, so a LIMIT message should be sent")
			}
		default:
			t.Error("Query is put in the queue, but the queue is over its limit, so a LIMIT message should be sent")
		}
	}
}
