package ns

import (
	"testing"

	"codeberg.org/miekg/dns"
)

func TestShouldRefuses(t *testing.T) {

	// invalid Name format
	REFtxtRR1 := dns.TXT{Hdr: dns.Header{Name: "test.com", Class: dns.ClassINET}}
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
	REFmsg2 := new(dns.Msg)
	REFmsg2.Question = append(REFmsg2.Question, &REFtxtRR2)
	REFmsg3 := new(dns.Msg)
	REFmsg3.Question = append(REFmsg3.Question, &REFtxtRR3)

	REFmsg4 := new(dns.Msg)
	REFmsg4.Question = append(REFmsg4.Question, &REFsoaRR)

	REFmsg5 := new(dns.Msg)
	REFmsg5.Question = append(REFmsg5.Question, &ACCtxtRR, &ACCtxtRR)
	ACCmsg6 := new(dns.Msg)
	ACCmsg6.Question = append(ACCmsg6.Question, &ACCtxtRR)

	if !shouldRefuse(REFmsg1) {
		t.Error("message 1 should be refused as it has invalid name format")

	}
	if !shouldRefuse(REFmsg2) {
		t.Error("message 2 should be refused as it has invalid name format")

	}
	if !shouldRefuse(REFmsg3) {
		t.Error("message 3 should be refused as it has invalid name format")

	}
	if !shouldRefuse(REFmsg4) {
		t.Error("message 4 should be refused as it has invalid type (soa)")

	}
	if !shouldRefuse(REFmsg5) {
		t.Error("message 5 should be refused as it has more than 1 question")

	}
	if shouldRefuse(ACCmsg6) {
		t.Error("message 6 should be accepted")
	}
}
