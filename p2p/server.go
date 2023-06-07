package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Settings struct {
	Name       string
	ListenAddr string
	Version    string
}

type Server struct {
	Settings
	mx              sync.RWMutex
	ln              net.Listener
	peers           map[net.Addr]*Peer
	addPeerCh       chan *Peer
	removePeerCh    chan net.Addr
	closeCh         chan struct{}
	closeAcceptLoop chan struct{}
	msgCh           chan Message
}

func NewServer(settings Settings) *Server {
	return &Server{
		Settings:        settings,
		peers:           make(map[net.Addr]*Peer),
		addPeerCh:       make(chan *Peer, 10),
		removePeerCh:    make(chan net.Addr),
		closeCh:         make(chan struct{}),
		closeAcceptLoop: make(chan struct{}),
		msgCh:           make(chan Message),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Settings.ListenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.channelHandler()
	go s.acceptLoop()

	<-s.closeCh
	log.Debugf("[%s] Shutting down start", s.Settings.Name)

	return nil
}

func (s *Server) Close() {
	s.closeCh <- struct{}{}
	s.closeAcceptLoop <- struct{}{}
	s.closeConnections()
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	err = s.connectPeer(conn)
	if err != nil {
		return err
	}

	conn.Write([]byte("Ping"))

	return nil
}

func (s *Server) greetPeer(conn net.Conn) (*Peer, error) {
	self := &Peer{
		Name: s.Settings.Name,
		Uuid: rand.Int63n(math.MaxInt64),
	}
	err := gob.NewEncoder(conn).Encode(self)
	if err != nil {
		return nil, fmt.Errorf("Connect peer failed (%s): %w", conn.RemoteAddr(), err)
	}

	var peer Peer
	err = gob.NewDecoder(conn).Decode(&peer)
	if err != nil {
		return nil, fmt.Errorf("Peer get information (%s): %w", conn.RemoteAddr(), err)
	}

	peer.conn = conn
	fmt.Printf("[%s] Peer %+v\n", s.Settings.Name, peer)
	return &peer, nil
}

func (s *Server) connectPeer(conn net.Conn) error {
	peer, err := s.greetPeer(conn)
	if err != nil {
		return err
	}

	s.addPeerCh <- peer

	go s.handleConn(conn)
	return nil
}

func (s *Server) GetPeer(addr net.Addr) (*Peer, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	peer, ok := s.peers[addr]
	return peer, ok
}

func (s *Server) closeConnections() {
	for _, p := range s.peers {
		err := p.conn.Close()
		if err != nil {
			p.Write([]byte("Closing connection"))
			log.Errorf("Conn close err: %s", err)
		}
	}
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.closeAcceptLoop:
				log.Debug("Accept Closing")
				return
			default:
				log.Fatalf("Accept: %s", err)
			}
		}
		err = s.connectPeer(conn)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		i, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Infof("[%s] Peer disconnected %s", s.Settings.Name, err)
				break
			}
			log.Fatal(err)
		}

		msg := string(buf[:i])

		s.msgCh <- Message{From: conn.RemoteAddr(), Payload: bytes.NewBuffer(buf[:i])}
		// TODO: Remove this block, testing only
		time.Sleep(time.Second)
		answer := "Pong"
		if msg == "Pong" {
			answer = "Ping"
		}
		conn.Write([]byte(answer))
		// --------------------------------------
	}
	s.removePeerCh <- conn.RemoteAddr()
}

func (s *Server) channelHandler() {
	for {
		select {
		case peer := <-s.addPeerCh:
			if _, exist := s.GetPeer(peer.conn.RemoteAddr()); exist {
				log.Errorf("[%s] Peer (%s) already exists", s.Settings.Name, peer.conn.RemoteAddr())
				peer.conn.Close()
			}
			s.mx.Lock()
			s.peers[peer.conn.RemoteAddr()] = peer
			s.mx.Unlock()
			log.Infof("[%s] New peer connected %s", s.Settings.Name, peer.conn.RemoteAddr())
			log.Debugf("[%s] Peers: %+v", s.Settings.Name, s.peers)
		case addr := <-s.removePeerCh:
			s.mx.Lock()
			delete(s.peers, addr)
			s.mx.Unlock()
			log.Infof("[%s] Removed peer: %s", s.Settings.Name, addr)
		case msg := <-s.msgCh:
			payload, _ := io.ReadAll(msg.Payload)
			log.Infof("[%s] Message from (%s): %s", s.Settings.Name, msg.From, payload)
		}

	}
}
