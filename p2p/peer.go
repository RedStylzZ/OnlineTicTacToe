package p2p

import "net"

type Peer struct {
	conn net.Conn
	Name string
	Uuid int64
}

func (p *Peer) Write(msg []byte) error {
	_, err := p.conn.Write(msg)
	return err
}
