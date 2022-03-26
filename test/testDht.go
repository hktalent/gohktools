package main

// "github.com/hktalent/gohktools/lib/iface/protocol/PwnDht"
import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	disc "github.com/libp2p/go-libp2p-discovery"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

func StringsToAddrs(addrStrings []string) (maddrs []multiaddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := multiaddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}
func main() {
	// idHex := "68747470733a2f2f646874346861636b65722e353170776e2e636f6d68747470733a2f2f646874346861636b65722e353170776e2e636f6d"
	// idBytes, err := hex.DecodeString(idHex)
	// if err != nil {
	// 	panic(err)
	// }
	// var id [20]byte
	// n := copy(id[:], idBytes)
	// if n == 20 {
	// }
	// s, err := dht.NewServer(&dht.ServerConfig{
	// 	NodeId: id,
	// 	Conn:   mustListen(":0"),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// defer s.Close()
	// https://stun4dht4hackers.51pwn.com
	// var myDht = PwnDht.New()
	// myDht.Start()
	// fmt.Println(myDht)

	chatProtocol := "/51pwn/p2p/dht/1.1.0"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mm, err := StringsToAddrs([]string{"/ip4/0.0.0.0/udp/3344"})
	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(libp2p.ListenAddrs(mm...))

	if err != nil {
		panic(err)
	}
	// host := dht.DefaultBootstrapPeers

	dht, err := kaddht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	routingDiscovery := disc.NewRoutingDiscovery(dht)
	disc.Advertise(ctx, routingDiscovery, string(chatProtocol))
	peers, err := disc.FindPeers(ctx, routingDiscovery, string(chatProtocol))
	if err != nil {
		panic(err)
	}
	for _, peer := range peers {
		fmt.Println(peer)
		// notifee.HandlePeerFound(peer)
	}

}
