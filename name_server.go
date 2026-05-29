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

	var hdr = &dns.Header{Name: r.Question[0].Header().Name + dom, Class: dns.ClassINET}
	r.Reset() // re-use r
	r.Response = true

	if a, b := ns.cache.Get(r.Question[0]); b == true {
		r.Answer = append(r.Answer, *a...)
		r.Pack()
		io.Copy(w, r)
		return
	}

	ns.query_queue.Push(r.Question[0], ip)

	r.Answer = append(r.Answer, &dns.HINFO{Hdr: *hdr, HINFO: rdata.HINFO{Cpu: "QUEUE", Os: "QUEUE"}})

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
