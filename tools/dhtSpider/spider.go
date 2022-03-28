package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/hktalent/dht"
)

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []file `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

var (
	resUrl = "http://127.0.0.1:9200/dht_index/_doc/"
	// len = 40
	myPeerId = hex.EncodeToString([]byte("https://ee.51pwn.com")[:20])
)

/*
save to es server
1、create index
PUT /dht_index HTTP/1.1
host:127.0.0.1:9200
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15
Connection: close
Content-Type: application/json;charset=UTF-8
Content-Length: 413

{
  "settings": {
   "analysis": {
     "analyzer": {
       "default": {
         "type": "custom",
         "tokenizer": "ik_max_word",
         "char_filter": [
            "html_strip"
          ]
       },
       "default_search": {
         "type": "custom",
         "tokenizer": "ik_max_word",
         "char_filter": [
            "html_strip"
          ]
      }
     }
   }
  }
}

2、settings
PUT /dht_index/_settings HTTP/1.1
host:127.0.0.1:9200
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15
Connection: close
Content-Type: application/json;charset=UTF-8
Content-Length: 291

{
  "index.mapping.total_fields.limit": 10000,
 "number_of_replicas" : 0,
"index.translog.durability": "async",
"index.blocks.read_only_allow_delete":"false",
    "index.translog.sync_interval": "5s",
    "index.translog.flush_threshold_size":"100m",
   "refresh_interval": "30s"

}
*/
func sendReq(data []byte, id string) {
	fmt.Println("start send to ", resUrl, " es "+id)
	// Post "77beaaf8081e4e45adb550194cc0f3a62ebb665f": unsupported protocol scheme ""
	req, err := http.NewRequest("POST", resUrl+id, bytes.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	// 取消全局复用连接
	// tr := http.Transport{DisableKeepAlives: true}
	// client := http.Client{Transport: &tr}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	// keep-alive
	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close() // resp 可能为 nil，不能读取 Body
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	// body, err := ioutil.ReadAll(resp.Body)
	// _, err = io.Copy(ioutil.Discard, resp.Body) // 手动丢弃读取完毕的数据
	// json.NewDecoder(resp.Body).Decode(&data)
	fmt.Println("[send request] " + id)
	// req.Body.Close()

	// go http.Post(resUrl, "application/json",, post_body)
}

/*
用来判断多少人使用该软件
可在这些使用者之间建立通讯
*/
func getMyPeer(d *dht.DHT) {
	fmt.Println("getMyPeer " + myPeerId)
	for {
		err := d.GetPeers(myPeerId)
		if err != nil && err != dht.ErrNotReady {
			fmt.Println(err)
		}

		if err == dht.ErrNotReady {
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}
}

func main() {
	resUrl := flag.String("resUrl", "", "Elasticsearch url, eg: http://127.0.0.1:9200/dht_index/_doc/")
	if "" == *resUrl {
		*resUrl = "http://127.0.0.1:9200/dht_index/_doc/"
	}
	flag.Parse()
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	nX := 10
	// blackListSize, requestQueueSize, workerQueueSize
	w := dht.NewWire(65536, 1024*nX, 256*nX)
	go func() {
		for resp := range w.Response() {
			metadata, err := dht.Decode(resp.MetadataInfo)
			if err != nil {
				continue
			}
			info := metadata.(map[string]interface{})

			if _, ok := info["name"]; !ok {
				continue
			}

			bt := bitTorrent{
				InfoHash: hex.EncodeToString(resp.InfoHash),
				Name:     info["name"].(string),
			}

			if v, ok := info["files"]; ok {
				files := v.([]interface{})
				bt.Files = make([]file, len(files))

				for i, item := range files {
					f := item.(map[string]interface{})
					bt.Files[i] = file{
						Path:   f["path"].([]interface{}),
						Length: f["length"].(int),
					}
				}
			} else if _, ok := info["length"]; ok {
				bt.Length = info["length"].(int)
			}

			data, err := json.Marshal(bt)
			if err == nil {
				if 0 < len(*resUrl) {
					sendReq(data, bt.InfoHash)
				}
				// fmt.Printf("%s\n\n", data)
			}
		}
	}()
	go w.Run()

	config := dht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(config)
	// d.Mode = &dht.newNode(myPeerId, "", config.Address)
	d.OnGetPeersResponse = func(infoHash string, peer *dht.Peer) {
		if infoHash == myPeerId {
			fmt.Printf("my private net: <%s:%d>\n", peer.IP, peer.Port)
		} else if 0 < len(*resUrl) {
			sendReq([]byte(fmt.Sprintf("{\"ip\":\"%s\",\"port\":%d,\"type\":\"peer\"}", peer.IP, peer.Port)), fmt.Sprintf("%s_%d", peer.IP, peer.Port))
			fmt.Printf("peer info : %s:%d\n", peer.IP, peer.Port)
		}
	}
	go getMyPeer(d)
	fmt.Println("start run ...")
	d.Run()
}