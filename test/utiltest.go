package main

import (
	"fmt"
	"net"
)

// 获取机器macdiz，构建全球唯一机器标识
func getMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

// 获取可用绑定端口
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
	s, _ := getMacAddr()
	fmt.Println(s)

}
