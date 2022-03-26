package protocol

import (
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
)

type PwnMDns struct {
	Config  map[string]string
	Notifee chan peer.AddrInfo
	ser     mdns.mdnsService
	host
}

func New() *PwnMDns {
	return &PwnMDns{}
}
func (n *PwnMDns) SetConfig(i interface{}) {

}
func (n *PwnMDns) GetConfig() interface{} {
	return n.Config
}
func (n *PwnMDns) Start() {

	n.Notifee.PeerChan = make(chan peer.AddrInfo)
	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", "0.0.0.0", 0))

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	n.host = host
	if err != nil {
		panic(err)
	}

	// An hour might be a long long period in practical applications. But this is fine for us
	n.ser = mdns.NewMdnsService(host, n.Config["serviceName"], n.Notifee)
	if err := n.ser.Start(); err != nil {
		panic(err)
	}
}
func (n *PwnMDns) Stop() {
	n.ser.Close()
}

func (n *PwnMDns) LogNoti(i interface{}) {

}
func (n *PwnMDns) DataNoti(i interface{}) {

}
