package main

import (
	"log"

	"github.com/insomniacslk/dhcp/dhcpv6/client6"
)

func RelayFourMsgTransaction(client *client6.Client, intf string) {
	// Exchange runs a Solicit-Advertise-Request-Reply transaction on the
	// specified network interface, and returns a list of DHCPv6 packets
	// (a "conversation") and an error if any. Notice that Exchange may
	// return a non-empty packet list even if there is an error. This is
	// intended, because the transaction may fail at any point, and we
	// still want to know what packets were exchanged until then.
	// A default Solicit packet will be used during the "conversation",
	// which can be manipulated by using modifiers.
	conversation, err := client.Exchange(intf)

	// Summary() prints a verbose representation of the exchanged packets.
	for _, packet := range conversation {
		log.Print(packet.Summary())
	}
	// error handling is done *after* printing, so we still print the
	// exchanged packets if any, as explained above.
	if err != nil {
		log.Fatal(err)
	}

}
