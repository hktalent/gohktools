package socks5

import (
	"fmt"
	"log"
	"os"

	nodeInfo "github.com/hktalent/gohktools/lib/node"
	utils "github.com/hktalent/gohktools/lib/utils"

	"github.com/armon/go-socks5"
)

type Socks5ServerInfo struct {
	Port int
	Ips  []string
}

func GetSocks5Config(key string) *socks5.Config {
	// Create a socks server
	creds := socks5.StaticCredentials{
		"51pwn": key,
	}
	cator := socks5.UserPassAuthenticator{Credentials: creds}
	return &socks5.Config{
		AuthMethods: []socks5.Authenticator{cator},
		// Resolver:    context.NameResolver{context.Background(), "8.8.8.8"},
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

/*
传入key
返回本机所有绑定的ip和端口信息给回调
*/
func Socks5Server(key string, fnCbk func(Socks5ServerInfo, interface{}), wg interface{}) error {
	nP, err := utils.GetFreePort("tcp")
	if err != nil {
		return err
	}
	conf := GetSocks5Config(key)

	// conf := &socks5.Config{}
	server, err := socks5.New(conf)
	nodeInfo.G_nodeInfo.Socks5Port = nP
	nodeInfo.G_nodeInfo.Socks5Server = server
	if err != nil {
		return err
	}
	go func() {
		// fmt.Println(fmt.Sprintf("in port is %v", nP))
		if err := server.ListenAndServe("tcp", fmt.Sprintf(":%v", nP)); err != nil {
			panic(err)
		}
	}()

	_, ips, err := utils.GetMacAddr()
	if nil != err {
		return err
	}
	nodeInfo.G_nodeInfo.Ips = ips
	ssi := Socks5ServerInfo{Port: nP, Ips: ips}
	go fnCbk(ssi, wg)
	return nil
}

// // github.com/txthinking/socks5
// func Socks4Client() {
// 	tcpTimeout := 100000
// 	udpTimeout := 100000
// 	c, _ := socks5.NewClient(server, username, password, tcpTimeout, udpTimeout)
// 	conn, _ := c.Dial(network, addr)
// }
