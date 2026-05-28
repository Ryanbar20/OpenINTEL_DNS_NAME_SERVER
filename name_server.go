package ns

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codeberg.org/miekg/dns"
)

type NameServer struct {
	query_queue Queue
	cache       Cache
}

func handle(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {

}

func serve(net string) {
	addr := fmt.Sprintf("[::]:%d", NAME_SERVER_PORT)
	server := &dns.Server{Addr: addr, Net: net, ReusePort: true, MaxTCPQueries: -1}
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to setup the "+net+" server: %s", err.Error())
	}
}

func newNameServer(cache_limit int, queue_limit int) *NameServer {
	return &NameServer{query_queue: *newQueue(queue_limit), cache: *newCache(cache_limit)}
}

func (ns *NameServer) start() {

	dns.HandleFunc(dom, handle)

	go query_thread(&ns.query_queue, &ns.cache)
	go serve("udp")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping", s)
}
