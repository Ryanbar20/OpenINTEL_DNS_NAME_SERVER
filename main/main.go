package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"strconv"

	ns "github.com/Ryanbar20/OpenINTEL_DNS_NAME_SERVER"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <cache_limit> <queue_limit> <ip:port>\n", os.Args[0])
		os.Exit(1)
	}

	cache_limit, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("cache limit (first argument) must be an integer")
	}

	queue_limit, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("queue limit (second argument) must be an integer")
	}

	addr := os.Args[3]
	var host string
	var portStr string
	if host, portStr, err = net.SplitHostPort(addr); err != nil {
		log.Fatal("invalid address")
	}

	ip, err := netip.ParseAddr(host)
	if err != nil {
		log.Fatal("invalid ip")
	}

	if ip.Is6() {
		host = "[" + host + "]" // brackets are needed for miekgdns
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("port must be an integer")
	}

	nameserver := ns.NewNameServer(cache_limit, queue_limit, host, port, ns.MEMORY_LIMIT)
	nameserver.Start()
}
