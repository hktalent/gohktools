package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	// "github.com/elastic/go-elasticsearch/v7"
	elastic "github.com/olivere/elastic/v7"
)

/*
// elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
./T3 -urlListsFile /Users/${USER}/MyWork/vulScanPro/T3.txt

PUT _scripts/upPswd HTTP/1.1
host:127.0.0.1:9200
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15
Connection: close
Content-Type: application/json;charset=UTF-8
Content-Length: 256

{
    "script": {
        "lang": "painless",
        "source": "if(ctx._source.containsKey('weblogic')){if(params.containsKey('T3')){ctx._source.weblogic['T3']=params['T3']}if(params.containsKey('GIOP')){ctx._source.weblogic['GIOP']=params['GIOP']}}"
    }
}

*/
var (
	urlListsFile, esUrl, numThread, tagName string
	saveIps                                 bool
	ctx11                                   context.Context
	fd                                      *os.File
	ins                                     []chan int
)

/*
域名转换为ip
*/
func addIps(url string) {
	ip := url
	if -1 < strings.Index(url, "://") {
		a := regexp.MustCompile(`://`)
		ip = a.Split(url, -1)[1]
	}
	if -1 < strings.Index(ip, "/") {
		a := regexp.MustCompile(`/`)
		ip = a.Split(ip, -1)[0]
	}
	if -1 == strings.Index(ip, ":") {
		ip += ":80"
	}
	domain := ip
	matched, _ := regexp.Match(`(\d{1,3}\.){3}\d{1,3}`, []byte(domain))
	if matched {
		return
	}
	// fmt.Println("[", ip, "]")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip) //创建一个TCPAddr
	if err != nil {
		fmt.Println(domain, err)
		return
	}

	ip = fmt.Sprintf("%s", tcpAddr.IP)
	fmt.Println(domain, " ", ip)
	fd_content := strings.Join([]string{domain, " ", ip, "\n"}, "")

	fd, _ = os.OpenFile("Ips.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd.Write([]byte(fd_content))
	fd.Close()
}

func IndexPrice(tweet map[string]interface{}) {

	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		fmt.Println("%v", err)
		return
	}

	bulkRequest := client.Bulk()
	request := elastic.NewBulkUpdateRequest()
	req := request.Index("51pwn_index").Type("_doc").Id(tweet["url"].(string)).Doc(tweet).DocAsUpsert(true)
	bulkRequest = bulkRequest.Add(req)
	bulkResponse, err := bulkRequest.Do(ctx11)
	if err != nil {
		fmt.Println(err)
		return
	}
	if bulkResponse != nil {
		fmt.Println(tweet["url"].(string))
		fmt.Println(bulkResponse)
	}

}

// 解决相同多次请求的问题
func sendReq(szRst, url1 string, szType, ip string) {
	type mi = map[string]interface{}
	// data1 := fmt.Sprintf(`{"url":"%s","ip":"%s","tagName":"%s","weblogic":{"%s":"%s"}}`, url1, ip, tagName, szType, url.QueryEscape(szRst))
	xxData := mi{"url": url1, "ip": ip, "tagName": tagName, "weblogic": mi{szType: url.QueryEscape(szRst)}}
	IndexPrice(xxData)
	// oSave := json.NewDecoder(`{"id":url,"url":url,"weblogic":{"T3":szRst}}`)

	// escapeUrl := url.QueryEscape(url1)
	// post_body := bytes.NewReader([]byte(fmt.Sprintf(`{"doc_as_upsert":"true","doc": {"url":"%s","ip":"%s","tagName":"%s","weblogic":{"%s":"%s"}}}`, url1, ip, tagName, szType, url.QueryEscape(szRst))))
	// fmt.Println(szType, " ", ip)
	// req, err := http.NewRequest("POST", esUrl+escapeUrl+"/_update", post_body)
	// if err == nil {
	// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
	// 	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	// 	req.Header.Set("Connection", "close")
	// 	// var rsp http.Response
	// 	rsp, err1 := http.DefaultClient.Do(req)
	// 	if err1 != nil {
	// 		fmt.Println(err1)
	// 		return
	// 	}
	// 	d, err := ioutil.ReadAll(rsp.Body)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	if -1 < strings.Index(string(d), `"successful":1`) {
	// 		fmt.Println(fmt.Sprintf("%s : %s  save ok", szType, url1))
	// 	} else {
	// 		fmt.Println(string(d))
	// 	}
	// 	defer rsp.Body.Close()
	// 	defer req.Body.Close()
	// }

	//  "http://127.0.0.1:9201"
	// addresses := []string{"http://127.0.0.1:9200"}
	// config := elasticsearch.Config{
	// 	Addresses: addresses,
	// 	Username:  "",
	// 	Password:  "",
	// 	CloudID:   "",
	// 	APIKey:    "",
	// }
	// es, _ := elasticsearch.NewClient(config)
	// // log.Println(elasticsearch.Version)
	// // log.Println(es.Info())
	// var buf bytes.Buffer
	// doc := map[string]interface{}{
	// 	"doc": map[string]interface{}{
	// 		"url":     url1,
	// 		"ip":      ip,
	// 		"tagName": tagName,
	// 		"weblogic": map[string]interface{}{
	// 			szType: url.QueryEscape(szRst),
	// 		},
	// 	},
	// }
	// if err := json.NewEncoder(&buf).Encode(doc); err != nil {
	// 	return
	// }
	// res, err := es.Update("51pwn_index", url1, &buf, es.Update.WithDocumentType("doc"))
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Panicln(res.Body)
	// defer res.Body.Close()
	// fmt.Println("sendReq end")
	// go http.Post(resUrl, "application/json",, post_body)
}

func senData(url string, szData, szCheck, szType string) {

	ip := url
	if -1 < strings.Index(url, "://") {
		a := regexp.MustCompile(`://`)
		ip = a.Split(url, -1)[1]
	}
	if -1 < strings.Index(ip, "/") {
		a := regexp.MustCompile(`/`)
		ip = a.Split(ip, -1)[0]
	}
	if -1 == strings.Index(ip, ":") {
		ip += ":80"
	}
	// fmt.Println(ip)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip) //创建一个TCPAddr
	if err != nil {
		return
	}

	tcpCoon, err := net.DialTCP("tcp4", nil, tcpAddr) //建立连接
	if err != nil {
		// fmt.Println("DialTCP:")
		// fmt.Println(err)
		return
	}
	defer tcpCoon.Close() //关闭
	sendData := szData
	n, err := tcpCoon.Write([]byte(sendData)) //发送数据
	if err != nil {
		// fmt.Println("Write:")
		// fmt.Println(err)
		return
	}
	recvData := make([]byte, 2048)
	n, err = tcpCoon.Read(recvData) //读取数据
	if err != nil {
		// fmt.Println("Read:")
		// fmt.Println(err)
		return
	}
	recvStr := string(recvData[:n])
	matched, _ := regexp.Match(szCheck, []byte(recvStr))

	if matched {
		// fmt.Println(recvStr)
		// fmt.Println(url)
		sendReq(recvStr, url, szType, fmt.Sprintf("%s", tcpAddr.IP))
	}
}

// if r and b'GIOP' in r and b'/bea_wls_internal/classes/' in r and b':weblogic/corba/cos/naming/NamingContextAny:' in r:
func sendGIOP(url string) {
	senData(url, "GIOP\x01\x02\x00\x03\x00\x00\x00\x17\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x0bNameService", "GIOP.*\\/bea_wls_internal\\/classes\\/.*:weblogic\\/corba\\/cos\\/naming\\/NamingContextAny:", "GIOP")
}

func sendT3(url string) {
	senData(url, "t3 12.2.1\nAS:255\nHL:19\nMS:10000000\n\n", "(HELO:)|(weblogic\\.security\\.net\\.FilterException)", "T3")
}

func main() {
	flag.StringVar(&urlListsFile, "urlListsFile", "/Users/"+os.Getenv("USER")+"/MyWork/vulScanPro/T3.txt", "url lists text file")
	flag.StringVar(&tagName, "tagName", "gov,yh", "set list tag name,eg: gov")

	flag.BoolVar(&saveIps, "saveIps", false, "save Ips to Ips.txt")

	flag.StringVar(&esUrl, "esUrl", "http://127.0.0.1:9200/51pwn_index/_doc/", "es server url")
	// flag.StringVar(&numThread, "numThread", "16", "threads number")
	flag.Parse()
	// Create a context object for the API calls
	ctx11 = context.Background()
	wg := &sync.WaitGroup{}

	file, err := os.OpenFile(urlListsFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		line = strings.TrimSpace(line)
		wg.Add(1)
		go func(line string) {
			if 7 < len(line) {
				if -1 == strings.Index(line, "://") {
					line = "http://" + line
				}
				if saveIps {
					addIps(line)
				} else {
					sendGIOP(line)
					sendT3(line)
				}
			}
			wg.Done()
		}(line)
	}
	wg.Wait()
}
