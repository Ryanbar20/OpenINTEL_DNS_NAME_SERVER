package ns

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"syscall"

	"codeberg.org/miekg/dns"
	"codeberg.org/miekg/dns/rdata"
)

type NameServer struct {
	query_queue Queue
	cache       Cache
}

func handle(ns *NameServer, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
	if err := r.Unpack(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	var hdr = &dns.Header{Name: r.Question[0].Header().Name, Class: dns.ClassINET}
	r.Reset() // re-use r
	r.Response = true

	// check if cache-hit
	if a, b := ns.cache.Get(r.Question[0]); b == true {
		fmt.Println("cache hit")
		fmt.Println(r.Question[0].String())
		r.Answer = append(r.Answer, *a...)
		r.Pack()
		io.Copy(w, r)
		return
	}
	fmt.Println("cache miss")
	fmt.Println(r.Question[0].String())
	fmt.Println(ns.cache.data)

	var ip netip.Addr
	switch a := w.RemoteAddr().(type) {
	case *net.UDPAddr:
		ip, _ = netip.AddrFromSlice(a.IP)
	case *net.TCPAddr:
		ip, _ = netip.AddrFromSlice(a.IP)
	}
	if ip.Is4In6() {
		ip = netip.AddrFrom4(ip.As4())
	}

	hdr.TTL = 0                                             // set TTL to 0 so that wait and limit messages are not cached
	if i := ns.query_queue.FindRR(r.Question[0]); i != -1 { // check if RR is in cache
		r.Answer = append(r.Answer, &dns.HINFO{Hdr: *hdr, HINFO: rdata.HINFO{Cpu: "WAIT", Os: fmt.Sprintf("You are in queue position %d", i)}})
	} else if i := ns.query_queue.FindIP(ip); i != -1 { // check if user is in cache
		r.Answer = append(r.Answer, &dns.HINFO{Hdr: *hdr, HINFO: rdata.HINFO{Cpu: "LIMIT", Os: fmt.Sprintf("You already have a query in queue position %d", i)}})
	} else { // push to queue
		queue_index := ns.query_queue.Push(r.Question[0], ip)
		if queue_index == -1 { // if queue limit was reached
			r.Answer = append(r.Answer, &dns.HINFO{Hdr: *hdr, HINFO: rdata.HINFO{Cpu: "LIMIT", Os: "Queue limit reached"}})
		} else { // if queue limit was not reached
			r.Answer = append(r.Answer, &dns.HINFO{
				Hdr:   *hdr,
				HINFO: rdata.HINFO{Cpu: "WAIT", Os: fmt.Sprintf("You are in queue position %d", queue_index)},
			})
		}
	}

	r.Pack()
	io.Copy(w, r)
}

func serve(net string) {
	addr := fmt.Sprintf("[::]:%d", NAME_SERVER_PORT)
	server := &dns.Server{Addr: addr, Net: net, ReusePort: true, MaxTCPQueries: -1}
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to setup the "+net+" server: %s", err.Error())
	}
}

func NewNameServer(cache_limit int, queue_limit int) *NameServer {
	return &NameServer{query_queue: *newQueue(queue_limit), cache: *newCache(cache_limit)}
}

func (ns *NameServer) Start() {

	dns.HandleFunc(dom, func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) { handle(ns, ctx, w, r) })

	go query_thread(&ns.query_queue, &ns.cache)
	go serve("udp")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping", s)
}
