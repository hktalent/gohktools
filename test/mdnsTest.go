package main

import (
	"fmt"

	mdns "github.com/hktalent/gohktools/lib/mdns"
)

// import (

// 	mdns "github.com/hktalent/gohktools/lib"
// )
func main() {
	name := "xxxx.xx"
	mdns.Server4mdns(name, "123.33.23:55")
	mdns.Client4mdns(name)
	fmt.Printf("start")
}
