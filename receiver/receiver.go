package main

import (
	"fmt"
	"log"
	"net"
)

const (
	broadcastAddr   = "224.0.0.1:9991"
	maxDatagramSize = 8192
)

func main() {
	// Connect to the multicast address.
	addr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		log.Fatal("failed to resolve UDP address: " + err.Error())
	}
	multicastListener, err := net.ListenMulticastUDP("udp", nil, addr)
	multicastListener.SetReadBuffer(maxDatagramSize)

	// Just wait for incoming packages.
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := multicastListener.ReadFromUDP(b)
		if err != nil {
			log.Fatal("could not read multicast package: " + err.Error())
		}
		fmt.Printf("read %d bytes from %s: '%s'\n", n, src.String(), string(b))
	}
}
