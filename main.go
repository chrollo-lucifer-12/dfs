package main

import (
	"fmt"
	"log"

	"github.com/chrollo-lucider-12/dfs/p2p"
)

func OnPeer(p2p.Peer) error {
	return nil
}

func main() {
	trOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
		OnPeer:        OnPeer,
		Logger:        p2p.NewSlogLogger(1),
	}

	tr := p2p.NewTCPTransport(trOpts)

	go func() {
		for msg := range tr.Consume() {
			fmt.Println(msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
