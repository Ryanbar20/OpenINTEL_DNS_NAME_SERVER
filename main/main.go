package main

import (
	"codeberg.org/miekg/dns"
)

var dom = "test.nl."
var hdr = &dns.Header{Name: dom, Class: dns.ClassINET}

func main() {

	// dns.HandleFunc(dom, handle_query)

	// go query_thread()
	// go serve("tcp")
	// go serve("udp")

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// s := <-sig
	// fmt.Printf("Signal (%s) received, stopping", s)

}
