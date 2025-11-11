package main

import (
	"bytes"
	"log"
	"log/slog"
	"time"

	"github.com/chrollo-lucider-12/dfs/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	logger := p2p.NewSlogLogger(slog.LevelInfo)
	tcpTransport := p2p.NewTCPTransport(p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Logger:        logger,
		Decoder:       p2p.NOPDecoder{},
	})
	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr[1:] + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(1 * time.Second)
	go s2.Start()
	time.Sleep(1 * time.Second)

	data := bytes.NewReader([]byte("sahil"))
	s2.StoreData("sahil", data)

	select {}
}
