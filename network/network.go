package network

import (
	"log"
	"net"
	"time"

	"github.com/jackpal/gateway"
	"github.com/tatsushid/go-fastping"
	"gitlab.com/MikeTTh/env"
)

// TODO: this could be made event driven

func HasNetwork() bool {
	addrs, err := gateway.DiscoverGateways()
	if err != nil {
		log.Println("Failed to discover gateways. Assuming no network...")
		return false
	}

	if env.Bool("FRIG_NO_PING", false) {
		// pinging does not work in all network scenarios,
		// so alternatively could alternatively consider the presence of the default gateway as a sign of being online
		return true
	}

	success := false

	p := fastping.NewPinger()
	_, _ = p.Network("udp")
	for _, addr := range addrs {
		err = p.AddIP(addr.String())
		if err != nil {
			log.Printf("Failed to add IP address %s for pinging: %v", addr.String(), err)
		}
	}
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		success = true
	}
	err = p.Run()
	if err != nil {
		log.Println("Failed to ping the default gateway:", err)
		success = false
	}
	if !success {
		log.Println("No response from pinging the default gateway")
	}

	return success
}
