package PwnDht

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/protocol"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/multiformats/go-multiaddr"
)

var logger = log.Logger("rendezvous")

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

type PwnDht struct {
	Config         map[string]string
	Notifee        chan peer.AddrInfo
	BootstrapPeers []multiaddr.Multiaddr
}

func New() *PwnDht {
	return &PwnDht{
		Config: map[string]string{
			"rendezvous": "info",
			"pid":        "/51pwn4dht/1.1.0",
			"peerID":     "xxxxxx",
			"bindIp":     "/ip4/0.0.0.0/tcp/0"},
		BootstrapPeers: dht.DefaultBootstrapPeers}
}

func handleStream(stream network.Stream) {
	logger.Info("Got a new stream!")
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func (n *PwnDht) SetConfig(i interface{}) {}
func (n *PwnDht) GetConfig() interface{} {
	return n.Config
}

func (n *PwnDht) StringsToAddrs(addrStrings []string) (maddrs []multiaddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := multiaddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}
func (n *PwnDht) Start() {
	mm, err := n.StringsToAddrs([]string{n.Config["bindIp"]})
	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(libp2p.ListenAddrs(mm...))
	if err != nil {
		panic(err)
	}
	host.SetStreamHandler(protocol.ID(n.Config["pid"]), handleStream)
	ctx := context.Background()
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}
	logger.Info("host ", host)
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range n.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				logger.Warning(err)
			} else {
				logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, n.Config["rendezvous"])
	peerChan, err := routingDiscovery.FindPeers(ctx, n.Config["rendezvous"])
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		logger.Debug("Found peer:", peer)

		logger.Debug("Connecting to:", peer)
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(n.Config["peerID"]))

		if err != nil {
			logger.Warning("Connection failed:", err)
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}

		logger.Info("Connected to:", peer)
	}

	select {}

}
func (n *PwnDht) Stop() {}

func (n *PwnDht) LogNoti(i interface{})  {}
func (n *PwnDht) DataNoti(i interface{}) {}

// var (
// 	G_PwnDht = New()
// )
