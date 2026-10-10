package p2p

type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(_ Peer) error {
	return nil
}
