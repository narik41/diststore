package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	// conn is the underlying  connection of the peer
	conn net.Conn
	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (t *TCPPeer) Close() error {
	return t.conn.Close()
}

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOps
	listener net.Listener
	message  chan Message
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
		message:         make(chan Message),
	}
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

func (t *TCPTransport) Consume() <-chan Message {
	return t.message
}

func (t *TCPTransport) startAcceptLoop() {
	fmt.Printf("TCP transport listening on %s\n", t.ListenAddr)
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error:%s \n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Println("TCP transport closing connection")
		conn.Close()
	}()
	peer := NewTCPPeer(conn, true)
	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP decode error:%s \n", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		t.message <- *msg
		fmt.Printf("TCP decode:%s\n", msg)
	}
}
