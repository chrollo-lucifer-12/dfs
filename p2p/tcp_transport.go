package p2p

import (
	"errors"
	"io"
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
	Logger        Logger
	OnPeer        func(Peer) error
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

func (t *TCPTransport) Close() error {
	return t.listener.Close()
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
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			t.Logger.Error("error accepting connection", "err")
		}

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		if err != nil && !errors.Is(err, io.EOF) {
			t.Logger.Error("closing peer %s due to error", "err", err, "address", conn.RemoteAddr())
		}
		conn.Close()
	}()

	peer := NewTCPPeer(conn, false)

	if err = t.HandshakeFunc(peer); err != nil {
		t.Logger.Error("TCP handshake error ", "err", err)
		conn.Close()
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	msg := Message{}
	for {
		err = t.Decoder.Decode(conn, &msg)
		if err != nil {
			if ok := errors.Is(err, io.EOF); ok {
				return
			}
			t.Logger.Error("TCP read Error ", "err", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		t.rpcch <- msg

	}
}
