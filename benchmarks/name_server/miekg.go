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

var hdr = &dns.Header{Name: "history.openintel.nl.", Class: dns.ClassINET}

func reflect(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
	if err := r.Unpack(); err != nil {
		log.Fatalf("%s", err.Error())
	}
	r.Reset() // re-use r
	r.Response = true

	txt1 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test1"}}}
	txt2 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test2"}}}
	txt3 := dns.TXT{Hdr: *hdr, TXT: rdata.TXT{Txt: []string{"Test3"}}}

	// support only these queries and return their answers
	switch r.Question[0].(type) {
	case *dns.A:
		r.Answer = append(r.Answer, &txt1)
	case *dns.AAAA:
		r.Answer = append(r.Answer, &txt2)
	case *dns.NS:
		r.Answer = append(r.Answer, &txt3)
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
	dns.HandleFunc("20201001.google.nu.history.openintel.nl.", reflect)
	go serve("udp")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping", s)
}
