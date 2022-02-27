package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
./T3 -urlListsFile /Users/${USER}/MyWork/vulScanPro/T3.txt

*/
var (
	urlListsFile, esUrl, numThread string
	saveIps                        bool
	ins                            []chan int
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
	if -1 < strings.Index(ip, ":") {
		a := regexp.MustCompile(`:`)
		ip = a.Split(ip, -1)[0]
	}

	// fmt.Println(ip)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip) //创建一个TCPAddr
	if err != nil {
		return
	}
	ip = fmt.Sprintf("%s", tcpAddr.IP)
	fd, _ := os.OpenFile("Ips.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fmt.Println(ip)
	fd_content := strings.Join([]string{ip, "\n"}, "")
	fd.Write([]byte(fd_content))
	fd.Close()
}

// 解决相同多次请求的问题
func sendReq(szRst, url1 string, in chan int, szType, ip string) {
	// oSave := json.NewDecoder(`{"id":url,"url":url,"weblogic":{"T3":szRst}}`)
	escapeUrl := url.QueryEscape(url1)
	post_body := bytes.NewReader([]byte(fmt.Sprintf(`{"url":"%s","ip":"%s","weblogic":{"%s":"%s"}}`, url1, ip, szType, url.QueryEscape(szRst))))
	fmt.Println(szType, " ", ip)
	req, err := http.NewRequest("POST", esUrl+escapeUrl, post_body)
	if err == nil {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		req.Header.Set("Connection", "close")
		// var rsp http.Response
		rsp, err1 := http.DefaultClient.Do(req)
		if err1 != nil {
			fmt.Println(err1)
			return
		}
		d, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
		}
		if -1 < strings.Index(string(d), `"successful":1`) {
			fmt.Println(fmt.Sprintf("%s : %s  save ok", szType, url1))
		}
		defer rsp.Body.Close()
		defer req.Body.Close()
	}
	// fmt.Println("sendReq end")
	// go http.Post(resUrl, "application/json",, post_body)
}

func senData(url string, in chan int, szData, szCheck, szType string) {
	func() {
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
			// fmt.Println("ResolveTCPAddr:")
			// fmt.Println(err)
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
			sendReq(recvStr, url, in, szType, fmt.Sprintf("%s", tcpAddr.IP))
		}
	}()
	in <- 1
}

// if r and b'GIOP' in r and b'/bea_wls_internal/classes/' in r and b':weblogic/corba/cos/naming/NamingContextAny:' in r:
func sendGIOP(url string, in chan int) {
	senData(url, in, "GIOP\x01\x02\x00\x03\x00\x00\x00\x17\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x0bNameService", "GIOP.*\\/bea_wls_internal\\/classes\\/.*:weblogic\\/corba\\/cos\\/naming\\/NamingContextAny:", "GIOP")
}

func sendT3(url string, in chan int) {
	senData(url, in, "t3 12.2.1\nAS:255\nHL:19\nMS:10000000\n\n", "(HELO:)|(weblogic\\.security\\.net\\.FilterException)", "T3")
}

func main() {
	flag.StringVar(&urlListsFile, "urlListsFile", "/Users/"+os.Getenv("USER")+"/MyWork/vulScanPro/T3.txt", "url lists text file")
	flag.BoolVar(&saveIps, "saveIps", false, "save Ips to Ips.txt")

	flag.StringVar(&esUrl, "esUrl", "http://127.0.0.1:9200/51pwn_index/_doc/", "es server url")
	flag.StringVar(&numThread, "numThread", "16", "threads number")
	flag.Parse()
	numT, _e := strconv.Atoi(numThread)
	if _e != nil {
		numT = 16
	}
	ins := make(chan int, numT)
	defer close(ins)
	file, err := os.OpenFile(urlListsFile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	nCnt := 0
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
		if 10 < len(line) {
			if -1 == strings.Index(line, "://") {
				line = "http://" + line
			}
			if saveIps {
				go addIps(line)
				nCnt++
			}

			go sendGIOP(line, ins)
			go sendT3(line, ins)
			nCnt += 2
		}
	}
	// fmt.Println(len(ins))
	// for _, ch1 := range ins {
	// 	<-ch1
	// 	// fmt.Println(n)
	// 	close(ch1)
	// }
	fmt.Println("总：", nCnt)
	// XXOver:
	for {
		// select {
		_, ok := <-ins
		nCnt--
		// fmt.Println("cur：", nCnt, "  ", ok)
		fmt.Print("+")
		if !ok || 0 >= nCnt {
			break
		}
		// }
	}
	// close(ins)
}
