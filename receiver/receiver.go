package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/ipv4"
)

const (
	broadcastAddr   = "224.0.0.1:9991"
	ipBroadcastAddr = "224.0.0.1:9992"
	maxDatagramSize = 8192
)

func main() {
	// IP multicast requires interface name.
	var itfName string
	flag.StringVar(&itfName, "itfname", "", "the name of the interface")
	flag.Parse()

	if itfName == "" {
		listenToUDPMulticast()
	} else {
		listenToIPMulticast(itfName)
	}
}

func listenToIPMulticast(itfName string) {
	itf, err := net.InterfaceByName(itfName)
	if err != nil {
		log.Fatal("failed to get interface " + itfName + ": " + err.Error())
	}

	// Connect to multicast address.
	addr, err := net.ResolveUDPAddr("udp", ipBroadcastAddr)
	if err != nil {
		log.Fatal("failed to resolve IP address: " + err.Error())
	}
	multicastCon, err := net.ListenPacket("udp", ipBroadcastAddr)
	if err != nil {
		log.Fatal("failed to establish IP connection: " + err.Error())
	}

	packetCon := ipv4.NewPacketConn(multicastCon)
	if err := packetCon.JoinGroup(itf, &net.UDPAddr{IP: addr.IP}); err != nil {
		log.Fatal("failed to join group: " + err.Error())
	}
	fmt.Println("joining multicast group ...")

	buf := make([]byte, 1500)
	for {
		n, _, _, err := packetCon.ReadFrom(buf)
		if err == nil {
			log.Printf("received IP multicast data: %s\n", buf[0:n])
		}
	}
}

func listenToUDPMulticast() {
	// Connect to the multicast address.
	addr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		log.Fatal("failed to resolve UDP address: " + err.Error())
	}
	multicastListener, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("failed to listen to UDP address: " + err.Error())
	}
	multicastListener.SetReadBuffer(maxDatagramSize)

	// Just wait for incoming packages.
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := multicastListener.ReadFromUDP(b)
		if err != nil {
			log.Fatal("could not read multicast package: " + err.Error())
		}
		fmt.Printf("received UDP multicast data; %d bytes from %s: '%s'\n", n, src.String(), string(b))
	}
}
