package main

import (
	"log"
	"log/slog"

	"github.com/chrollo-lucider-12/dfs/p2p"
)

func OnPeer(p2p.Peer) error {
	return nil
}

func main() {
	logger := p2p.NewSlogLogger(slog.LevelInfo)
	tcpTransport := p2p.NewTCPTransport(p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Logger:        logger,
		Decoder:       p2p.NOPDecoder{},
	})
	fileServerOpts := FileServerOpts{
		StorageRoot:       "chrollo",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	s := NewFileServer(fileServerOpts)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
