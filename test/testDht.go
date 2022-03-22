package main

import (
	"fmt"

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

	fmt.Println(GenCid("https://stun4dht4hackers.51pwn.com"))

}
