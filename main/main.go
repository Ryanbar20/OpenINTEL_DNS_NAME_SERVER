package main

import (
	ns "github.com/Ryanbar20/OpenINTEL_DNS_NAME_SERVER"
)

func main() {

	nameserver := ns.NewNameServer(10, 10)

	nameserver.Start()
	// dns.HandleFunc(dom, handle_query)

	// go query_thread()
	// go serve("tcp")
	// go serve("udp")

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// s := <-sig
	// fmt.Printf("Signal (%s) received, stopping", s)

}
