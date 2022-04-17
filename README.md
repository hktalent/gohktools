[![Tweet](https://img.shields.io/twitter/url/http/Hktalent3135773.svg?style=social)](https://twitter.com/intent/follow?screen_name=Hktalent3135773) [![Follow on Twitter](https://img.shields.io/twitter/follow/Hktalent3135773.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=Hktalent3135773) [![GitHub Followers](https://img.shields.io/github/followers/hktalent.svg?style=social&label=Follow)](https://github.com/hktalent/)
[![Top Langs](https://profile-counter.glitch.me/hktalent/count.svg)](https://51pwn.com)

# gohktools
golang hack tools

## weblogic T3/GIOP batch testing，default result save to elasticsearch 7.x http://127.0.0.1:9200/51pwn_index/_doc/
```bash
./T3 -h
-esUrl "http://your.es.ipAndPort/51pwn_index/_doc/"
```

## how test T3
```bash
go generate -x
# update to laest
go mod download github.com/hktalent/dht
go mod tidy

go build T3.go

./T3 -urlListsFile /Users/${USER}/MyWork/vulScanPro/T3.txt

cat Ips.txt|sort -u|uniq>Ips1.txt; mv Ips1.txt Ips.txt
masscan --rate=5000 --top-ports 10000 -oX out.xml -iL Ips.txt

```

# tools
- DHT spider
https://github.com/hktalent/dht/blob/master/sample/spider/spider.go

- Social worker data[;:], stored in elasticsearch
tools/indexer.go

```bash
find  $HOME/sgk1/BreachCompilation/data|xargs -I % ./indexer -filename="%"
```


## others
```
go install github.com/anacrolix/torrent/cmd/...@latest

torrent download 'magnet:?xt=urn:btih:KRWPCX3SJUM4IMM4YF5RPHL6ANPYTQPU'
torrent metainfo testdata/debian-10.8.0-amd64-netinst.iso.torrent magnet
# https://github.com/anacrolix/torrent

```


<!--

import "github.com/olivere/elastic/v7"

elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
-->
