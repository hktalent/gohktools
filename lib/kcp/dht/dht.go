package dht

import (
	"context"
	"encoding/hex"
	"net"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/anacrolix/envpprof"
	"github.com/anacrolix/log"

	"github.com/anacrolix/dht/v2"
)

var (
	flags = struct {
		TableFile   string `help:"name of file for storing node info"`
		Addr        string `help:"local UDP address"`
		NoBootstrap bool
	}{
		Addr: ":0",
	}
	s *dht.Server
)

func loadTable() (err error) {
	added, err := s.AddNodesFromFile(flags.TableFile)
	log.Printf("loaded %d nodes from table file", added)
	return
}

func saveTable() error {
	return dht.WriteNodesToFile(s.Nodes(), flags.TableFile)
}

// func main() {
// 	stdLog.SetFlags(stdLog.LstdFlags | stdLog.Lshortfile)
// 	err := mainErr()
// 	if err != nil {
// 		log.Printf("error in main: %v", err)
// 		os.Exit(1)
// 	}
// }

// "https://dht4hacker.51pwn.com".encode('utf-8').hex()
func getSameID() []byte {
	idHex := "68747470733a2f2f646874346861636b65722e353170776e2e636f6d68747470733a2f2f646874346861636b65722e353170776e2e636f6d"
	idBytes, err := hex.DecodeString(idHex)
	if err != nil {
		panic(err)
	}
	var id [20]byte
	copy(id[:], idBytes)
	return id
}

/*
"/ip4/0.0.0.0/udp/9000/quic"
*/
func mustListen(addr string) net.PacketConn {
	ret, err := net.ListenPacket("udp", addr)
	if err != nil {
		panic(err)
	}
	return ret
}

func mainErr() error {
	conn := mustListen(":0")
	defer conn.Close()
	cfg := dht.NewDefaultServerConfig()
	cfg.Conn = conn
	cfg.Logger = log.Default.FilterLevel(log.Info)
	cfg.NoSecurity = false
	cfg.NodeId = getSameID()
	s, err = dht.NewServer(cfg)
	if err != nil {
		return err
	}
	http.HandleFunc("/debug/dht", func(w http.ResponseWriter, r *http.Request) {
		s.WriteStatus(w)
	})
	if flags.TableFile != "" {
		err = loadTable()
		if err != nil {
			return err
		}
	}
	log.Printf("dht server on %s, ID is %x", s.Addr(), s.ID())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		log.Printf("got signal: %v", <-ch)
		cancel()
	}()
	if !flags.NoBootstrap {
		go func() {
			if tried, err := s.Bootstrap(); err != nil {
				log.Printf("error bootstrapping: %s", err)
			} else {
				log.Printf("finished bootstrapping: %#v", tried)
			}
		}()
	}
	<-ctx.Done()
	s.Close()

	if flags.TableFile != "" {
		if err := saveTable(); err != nil {
			log.Printf("error saving node table: %s", err)
		}
	}
	return nil
}
