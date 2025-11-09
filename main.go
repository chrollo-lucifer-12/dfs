package main

import (
	"log"

	"github.com/chrollo-lucider-12/dfs/p2p"
)

func main() {
	listenAddr := ":3000"
	tr := p2p.NewTCPTransport(listenAddr)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
