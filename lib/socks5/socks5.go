package socks5

import (
	"fmt"
	"log"
	"os"
	"sync"

	utils "github.com/hktalent/gohktools/lib/utils"

	"github.com/armon/go-socks5"
)

type Socks5ServerInfo struct {
	Port int
	Ips  []string
}

/*
传入key
返回本机所有绑定的ip和端口信息给回调
*/
func Socks5Server(key string, fnCbk func(Socks5ServerInfo, sync.WaitGroup), wg sync.WaitGroup) error {
	nP, err := utils.GetFreePort("tcp")
	if err != nil {
		return err
	}
	// Create a socks server
	creds := socks5.StaticCredentials{
		"51pwn": key,
	}
	cator := socks5.UserPassAuthenticator{Credentials: creds}
	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{cator},
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
	}

	// conf := &socks5.Config{}
	server, err := socks5.New(conf)
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
	ssi := Socks5ServerInfo{Port: nP, Ips: ips}
	go fnCbk(ssi, wg)
	return nil
}
