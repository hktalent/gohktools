package utils

import (
	"fmt"
	"net"

	nanoid "github.com/aidarkhanov/nanoid/v2"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

/*
Creating a CID from scratch
*/
func GenCid(s string) string {
	pref := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}

	// And then feed it some data
	// bafkreico4cwnommwxfopkyny63sznyup24g2fzbry3kvmbrvfmkcgud35a
	c, err := pref.Sum([]byte(s))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println("Created CID: ", c)
	return fmt.Sprint(c)
}

func GetUuid() string {
	id, err := nanoid.New() //> "i25_rX9zwDdDn7Sg-ZoaH"
	if err != nil {
		return ""
	}
	return id
}

// 获取机器mac、ip地址，构建全球唯一机器标识
func GetMacAddr() ([]string, []string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	var as, ipsAll []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, nil, err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// fmt.Println("Current IP address : ", ipnet.IP.String())
				ipsAll = append(ipsAll, ipnet.IP.String())
				// fmt.Println(ipnet.IP.String())
			}
		}
	}

	return as, ipsAll, nil
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

// func main() {
// 	n, err := GetFreePort("udp")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("udp port is ", n)
// 	n, err = GetFreePort("tcp")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("tcp port is ", n)
// 	s, _, _ := getMacAddr()
// 	fmt.Println(s)

// }
