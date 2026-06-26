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
	query_queue  Queue
	cache        Cache
	ip           string
	port         int
	memory_limit string
}

// checks if a question should be refused. MiekgDNS automatically refuses multi-question queries but not 0 question queries
func shouldRefuse(r *dns.Msg) bool {
	if len(r.Question) != 1 {
		return true // no question
	}
	_, b := parseName(r.Question[0].Header().Name)
	if !b {
		return true // invalid format
	}
	switch r.Question[0].(type) {
	case *dns.A, *dns.AAAA, *dns.NS, *dns.MX, *dns.TXT:
		return false
	default:
		return true // unsupported type
	}
}

// handles a DNS query
func (ns *NameServer) handle(_ context.Context, w dns.ResponseWriter, r *dns.Msg) {
	if err := r.Unpack(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	r.Reset() // re-use r
	r.Response = true

	if shouldRefuse(r) {
		r.Rcode = dns.RcodeRefused
		r.Pack()
		io.Copy(w, r)
		return
	}

	var hdr = &dns.Header{Name: r.Question[0].Header().Name, Class: dns.ClassINET}
	// check if cache-hit
	if a, b := ns.cache.Get(r.Question[0]); b == true {
		r.Answer = append(r.Answer, *a...)
		r.Pack()
		io.Copy(w, r)
		return
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

	// set TTL to 0 so that wait and limit messages are not cached
	hdr.TTL = 0
	if i := ns.query_queue.FindRR(r.Question[0]); i != -1 { // check if RR is in queue
		r.Answer = append(r.Answer, &dns.HINFO{Hdr: *hdr, HINFO: rdata.HINFO{Cpu: "WAIT", Os: fmt.Sprintf("You are in queue position %d", i)}})
	} else if i := ns.query_queue.FindIP(ip); i != -1 { // check if user is in queue
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

// starts the name server main thread
func serve(net string, ip string, port int) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	server := &dns.Server{Addr: addr, Net: net, ReusePort: true, MaxTCPQueries: -1}
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to setup the "+net+" server: %s", err.Error())
	}
}

// create a new name server with a cache and queue limit, an ip and port and a memory limit (for DuckDB).
func NewNameServer(cache_limit int, queue_limit int, ip string, port int, memory_limit string) *NameServer {
	return &NameServer{query_queue: *newQueue(queue_limit), cache: *newCache(cache_limit), ip: ip, port: port, memory_limit: memory_limit}
}

// starts the name server
func (ns *NameServer) Start() {

	dns.HandleFunc(dom, ns.handle)

	go query_thread(&ns.query_queue, &ns.cache, ns.memory_limit)
	go serve("udp", ns.ip, ns.port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping", s)
}
