package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

const (
	broadcastAddr = "224.0.0.1:9991"
)

func main() {
	// Connect to multicast address.
	addr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		log.Fatal("failed to resolve UDP address: " + err.Error())
	}
	multicastCon, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("failed to establish UDP connection" + err.Error())
	}

	// Just keep sending out packages.
	i := 0
	for {
		multicastCon.Write([]byte("package " + strconv.Itoa(i)))
		time.Sleep(1 * time.Second)
		i++
	}
}
