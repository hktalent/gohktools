# gohktools
golang hack tools

## default result save to elasticsearch 7.x
```
http://127.0.0.1:9200/51pwn_index/_doc/
./T3 -h
-esUrl "http://your.es.ipAndPort/51pwn_index/_doc/"

```

## how test T3
```bash

go build T3.go

./T3 -urlListsFile /Users/${USER}/MyWork/vulScanPro/T3.txt

cat Ips.txt|sort -u|uniq>Ips1.txt; mv Ips1.txt Ips.txt
masscan --rate=5000 --top-ports 10000 -oX out.xml -iL Ips.txt

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
