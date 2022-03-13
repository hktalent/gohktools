package main

import (
	"fmt"
	"github.com/hktalent/gohktools/lib/mdns"
)

// import (

// 	mdns "github.com/hktalent/gohktools/lib"
// )
func main() {
	mdns.mdns.server4mdns("xxxx.xx", "123.33.23:55")
	fmt.Printf("start")
}
