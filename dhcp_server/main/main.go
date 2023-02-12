package main

import (
	"log"
	"net"
	"time"

	uuid "github.com/google/UUID"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
)

var (
	//SVR_UUID = os.Getenv("SVR_UUID")
	SVR_UUID = "b70ee641-c51f-4ed6-b34b-e060dfceff46"
	DUID     dhcpv6.Duid
)

func handler(conn net.PacketConn, peer net.Addr, r dhcpv6.DHCPv6) {
	log.Print(r.Summary())
	clntMsg, err := r.GetInnerMessage()
	if err != nil {
		log.Fatal(err)
	}

	var resp dhcpv6.DHCPv6
	switch clntMsg.Type() {

	// SOLICIT
	case dhcpv6.MessageTypeSolicit:
		if resp, err = dhcpv6.NewAdvertiseFromSolicit(clntMsg, dhcpv6.WithServerID(DUID), dhcpv6.WithIANA()); err != nil {
			log.Fatal(err)
		}
		log.Print(resp.Summary())

	// REQUEST
	case dhcpv6.MessageTypeRequest:
		ops := dhcpv6.Options{}
		ops.Add(dhcpv6.OptBootFileURL("1.2.3.4/file/boot"))
		if resp, err = dhcpv6.NewReplyFromMessage(clntMsg, dhcpv6.WithIANA(dhcpv6.OptIAAddress{IPv6Addr: net.ParseIP("ABCD::FFFF"), ValidLifetime: time.Second * 1000, Options: dhcpv6.AddressOptions{ops}})); err != nil {
			log.Fatal(err)
		}
		log.Print(resp.Summary())

	}
	if r.IsRelay() {

		rr, err := dhcpv6.NewRelayReplFromRelayForw(r.(*dhcpv6.RelayMessage), resp.(*dhcpv6.Message))
		if err != nil {
			log.Fatal(err)
		}
		resp = rr
	}
	if _, err := conn.WriteTo(resp.ToBytes(), peer); err != nil {
		log.Printf("failed to send resp: %v", err)
	} else {
		log.Print("response sent")
	}
}

func main() {

	if uuid, err := uuid.Parse(SVR_UUID); err != nil {
		log.Fatal(err)
	} else if id, err := uuid.MarshalBinary(); err != nil {
		log.Fatal(err)
	} else {
		DUID = dhcpv6.Duid{Type: dhcpv6.DUID_UUID, Uuid: id}
	}

	laddr := net.UDPAddr{
		IP:   net.ParseIP("::"),
		Port: dhcpv6.DefaultServerPort,
	}
	server, err := server6.NewServer("eth1", &laddr, handler, server6.WithSummaryLogger())
	if err != nil {
		log.Fatal(err)
	}

	server.Serve()
}
