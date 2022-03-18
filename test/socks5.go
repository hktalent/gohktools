package main

import (
	"fmt"
	"sync"

	socks5 "github.com/hktalent/gohktools/lib/socks5"
	utils "github.com/hktalent/gohktools/lib/utils"
)

func fnCbk(rst socks5.Socks5ServerInfo, wg interface{}) {
	wg1, ok := wg.(sync.WaitGroup)
	if ok {
		defer wg1.Done()
	}
	fmt.Println("port is ", rst.Port)
	// fmt.Println(fmt.Sprintf("port is %v", x1))
	a := rst.Ips
	for i, x := range a {
		fmt.Println(i, " ip ", x)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	err := socks5.Socks5Server(utils.GetUuid(), fnCbk, wg)
	if err != nil {
		fmt.Println(err)
	}
	wg.Wait()
	// time.Sleep(5 * time.Second)
	// server
	// fmt.Println(fmt.Sprintf(":%v", nP))
}
