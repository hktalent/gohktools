package mdns

import (
	"fmt"
	"os"

	"github.com/hashicorp/mdns"
)

func Server4mdns(name, ipInfo string) {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{ipInfo}
	service, _ := mdns.NewMDNSService(host, name, "", "", 8000, nil, info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()
}

func Client4mdns(name string) {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	mdns.Lookup(name, entriesCh)
	close(entriesCh)
}
