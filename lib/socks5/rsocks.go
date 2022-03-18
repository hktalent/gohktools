package socks5

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/armon/go-socks5"

	"github.com/hashicorp/yamux"
)

var session *yamux.Session

/*
反向代理：
1、A 可以连接 B，B无法连接A
2、A建立 socks5 server，并与A见你
返回本机所有绑定的ip和端口信息给回调
*/
func RSocks5(key, address string, nP int) error {
	conf := GetSocks5Config(key)
	server, err := socks5.New(conf)
	if err != nil {
		return err
	}
	var conn net.Conn
	log.Println("Connecting to far end")
	conn, err = net.Dial("tcp", address)
	if err != nil {
		return err
	}

	log.Println("Starting server")
	session, err = yamux.Server(conn, nil)
	if err != nil {
		return err
	}

	for {
		stream, err := session.Accept()
		log.Println("Acceping stream")
		if err != nil {
			return err
		}
		log.Println("Passing off to socks5")
		go func() {
			err = server.ServeConn(stream)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

// Catches yamux connecting to us
func ListenForSocks(address string) {
	log.Println("Listening for the far end")
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		log.Println("Got a client")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Errors accepting!")
		}
		// Add connection to yamux
		session, err = yamux.Client(conn, nil)
	}
}

// Catches clients and connects to yamux
func ListenForClients(address string) error {
	log.Println("Waiting for clients")
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		// TODO dial socks5 through yamux and connect to conn

		if session == nil {
			conn.Close()
			continue
		}
		log.Println("Got a client")

		log.Println("Opening a stream")
		stream, err := session.Open()
		if err != nil {
			return err
		}

		// connect both of conn and stream

		go func() {
			log.Println("Starting to copy conn to stream")
			io.Copy(conn, stream)
			conn.Close()
		}()
		go func() {
			log.Println("Starting to copy stream to conn")
			io.Copy(stream, conn)
			stream.Close()
			log.Println("Done copying stream to conn")
		}()
	}
}
