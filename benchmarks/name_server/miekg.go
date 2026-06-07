package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

var answerMap = make(map[string]dns.RR, 0)

func handle(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
	if err := r.Unpack(); err != nil {
		log.Fatalf("%s", err.Error())
	}
	r.Reset() // re-use r
	r.Response = true
	rr := answerMap[r.Question[0].String()]
	if rr != nil {
		r.Answer = append(r.Answer, rr)
	}

	r.Pack()
	io.Copy(w, r)
}

func serve(net string) {
	addr := "[::]:10000"
	server := &dns.Server{Addr: addr, Net: net, ReusePort: true, MaxTCPQueries: -1}
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to setup the "+net+" server: %s", err.Error())
	}
}

func main() {

	var hdr = &dns.Header{Name: "20201001.google.nu.history.openintel.nl.", Class: dns.ClassINET}

	// support only the following queries and return their answers
	// this corresponds to the cached answers in the pipeline
	q1 := dns.A{Hdr: *hdr}
	q2 := dns.AAAA{Hdr: *hdr}
	q3 := dns.NS{Hdr: *hdr}
	txt1 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test1"}}}
	txt2 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test2"}}}
	txt3 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test3"}}}
	answerMap[q1.String()] = &txt1
	answerMap[q2.String()] = &txt2
	answerMap[q3.String()] = &txt3

	dns.HandleFunc(hdr.Name, handle)
	go serve("udp")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping", s)
}
