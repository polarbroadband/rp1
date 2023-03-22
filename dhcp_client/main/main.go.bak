package main

import (
	"flag"
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6/client6"
)

var (
	iface = flag.String("i", "eth1", "Interface to configure via DHCPv6")
)

func main() {
	flag.Parse()
	log.Printf("Starting DHCPv6 client on interface %s", *iface)

	// NewClient sets up a new DHCPv6 client with default values
	// for read and write timeouts, for destination address and listening
	// address
	client := client6.NewClient()

	client.SimulateRelay = true
	client.RemoteAddr = net.Addr(&net.UDPAddr{
		IP:   net.ParseIP("fd00:1:f004::f:3"),
		Port: 547,
		Zone: "eth1",
	})

	FourMsgTransaction(client, *iface)
}
