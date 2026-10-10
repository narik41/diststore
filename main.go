package main

import (
	"distributed-file-storage/internal/p2p"
	"log"
)

func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
