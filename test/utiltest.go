package main

import (
	"fmt"
	"net"
)

func GetFreePort(szType string) (int, error) {
	// var addr, err
	if "tcp" == szType {
		addr, err := net.ResolveTCPAddr(szType, ":0")
		if err != nil {
			return 0, err
		}
		l, err := net.ListenTCP(szType, addr)
		defer l.Close()
		return l.Addr().(*net.TCPAddr).Port, nil
	} else {
		addr, err := net.ResolveUDPAddr(szType, ":0")
		if err != nil {
			return 0, err
		}
		l, err := net.ListenUDP(szType, addr)
		if err != nil {
			return 0, err
		}
		defer l.Close()

		return l.LocalAddr().(*net.UDPAddr).Port, nil
	}
}

func main() {
	n, err := GetFreePort("udp")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("udp port is ", n)
	n, err = GetFreePort("tcp")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("tcp port is ", n)

}
