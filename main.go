package main

import (
	"fmt"
	"log"

	"github.com/chrollo-lucider-12/dfs/p2p"
)

func main() {
	trOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
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
