package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	udpBroadcastAddr = "224.0.0.1:9991"
	ipBroadcastAddr  = "224.0.0.1:9992"
)

func main() {
	// IP multicast requires interface name.
	var itfName string
	flag.StringVar(&itfName, "itfname", "eth0", "the name of the interface")
	flag.Parse()

	udpMulticaster := createUDPMulticaster()
	ipMulticaster := createIPMulticaster(itfName)
	defer ipMulticaster.Close()

	// Just keep sending out packages.
	i := 0
	for {
		udpMulticaster.Write([]byte("UDP package " + strconv.Itoa(i)))
		ipMulticaster.Write([]byte("IP package " + strconv.Itoa(i)))

		time.Sleep(1 * time.Second)
		i++
	}
}

func createUDPMulticaster() *net.UDPConn {
	// Connect to multicast address.
	addr, err := net.ResolveUDPAddr("udp", udpBroadcastAddr)
	if err != nil {
		log.Fatal("failed to resolve UDP address: " + err.Error())
	}
	multicastCon, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("failed to establish UDP connection" + err.Error())
	}

	return multicastCon
}

type ipMulticaster struct {
	addr *net.UDPAddr
	con  net.PacketConn
}

func createIPMulticaster(itfName string) *ipMulticaster {
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

	return &ipMulticaster{
		addr: addr,
		con:  multicastCon,
	}
}

func (multicaster *ipMulticaster) Write(b []byte) {
	multicaster.con.WriteTo(b, multicaster.addr)
}

func (multicaster *ipMulticaster) Close() error {
	return multicaster.con.Close()
}
