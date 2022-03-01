package main

import (
	"context"
	"fmt"
	"net"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func main() {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		fmt.Print("DiscoverGateway ", err)
		return
	}
	// fmt.Print("%+v\n", gatewayIP)
	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		fmt.Print("GetExternalAddress ", err)
		return
	}
	ctx = context.Background()
    timeout time.Duration
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	conn, err := (&net.Dialer{}).DialContext(timeoutCtx, "udp", net.JoinHostPort(ip.String(), "5351"))
	fmt.Printf("External IP address: %v\n", response.ExternalIPAddress)
}
