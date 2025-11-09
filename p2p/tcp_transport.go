package p2p

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
func (p *TCPPeer) Close() error {
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	rpcch    chan Message
	listener net.Listener
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan Message),
	}
}

func (t *TCPTransport) Consume() <-chan Message {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error : %s\n", err)
		}

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, false)

	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP handshake error : %s\n", err)
		conn.Close()
		return
	}

	msg := Message{}
	for {
		if err := t.Decoder.Decode(conn, &msg); err != nil {
			fmt.Printf("TCP Error : %s\n", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		t.rpcch <- msg
		fmt.Printf("message: %+v\n", msg)
	}
}
