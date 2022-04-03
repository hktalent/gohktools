package main

import (
	"fmt"

	"github.com/hktalent/dht"
)

func main() {
	downloader := dht.NewWire(65535, 1024, 256)
	go func() {
		// once we got the request result
		for resp := range downloader.Response() {
			fmt.Println(resp.InfoHash, resp.MetadataInfo)
		}
	}()
	go downloader.Run()

	config := dht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		// request to download the metadata info
		downloader.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(config)

	d.Run()
}
