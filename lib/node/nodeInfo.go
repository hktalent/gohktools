package node

import (
	"github.com/armon/go-socks5"
)

// 	"fmt"
// 	"log"
// 	"os"

// 	utils "github.com/hktalent/gohktools/lib/utils"

type NodeInfo struct {
	Socks5Port int
	HttpsPort  int
	Ips        []string

	Socks5Server *socks5.Server
}

var G_nodeInfo = NodeInfo{}
