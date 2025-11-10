package main

import (
	"log"
	"log/slog"

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
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)
	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()
}
