package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"sync"
	"time"
)

const (
	ListenPort                = 20111
	CommandPort               = 20112
	QuietQueryType PacketType = 1
)

var QuietHours = []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

type Server struct {
	closer io.Closer
	Addr   string
	Port   int
	wg     sync.WaitGroup
	done   chan any
}

func NewServer(addr string, port int) *Server {
	return &Server{
		Addr: addr,
		Port: port,
		wg:   sync.WaitGroup{},
	}
}

func (s *Server) Shutdown() { close(s.done) }

func handleQuery(p Preamble, data []byte) ([]byte, error) {
	var q QuietQuery
	if q.Type() != p.PacketType {
		return nil, fmt.Errorf("packet type mismatch")
	}
	q.Preamble = &p

	if _, err := q.Unmarshal(data); err != nil {
		return nil, err
	}
	now := time.Now()

	log.Printf("Query: %v", q)

	idx := now.Hour()

	inQuietHours := QuietHours[idx] != 0

	wakeUpIn := 0
	for i := idx + 1; i <= len(QuietHours) && QuietHours[i-1] == 1; i++ {
		wakeUpIn += 1
		if wakeUpIn >= 24 {
			break
		}
	}

	resp := QuietReponse{
		WakeUpHour:  uint16(wakeUpIn),
		IsQuietTime: inQuietHours,
	}

	log.Printf("%+v", resp)

	return resp.Marshal()
}

type PacketType uint16

func (s *Server) listenForMessages(conn net.PacketConn) {
	process := func(conn net.PacketConn, packet []byte, addr net.Addr) {
		log.Printf("[INFO] processing packet %d", len(packet))
		p := &Preamble{}
		n, err := p.Unmarshal(packet)
		if err != nil {
			log.Printf("[ERR] error processing preamble: %s", err)
		}

		switch p.PacketType {
		case QuietQueryType:
			if data, err := handleQuery(*p, packet[n:]); err != nil {
				log.Printf("[ERR] error processing query: %s", err)
			} else {
				conn.WriteTo(data, addr)
			}
		default:
			log.Printf("[WARN] unknown packet type: %d", p.PacketType)
		}
	}
	select {
	case <-s.done:
		{
			s.wg.Done()
		}
	default:
		{
			for {
				var packet [512]byte
				n, addr, err := conn.ReadFrom(packet[:])
				if n >= 0 && n <= len(packet) {
					go process(conn, packet[:], addr)
				} else if err != nil {
					log.Printf("[ERR] error processing packet: %s", err)
				}
			}
		}
	}
}

func (s *Server) Listen() {
	addr := net.UDPAddrFromAddrPort(netip.MustParseAddrPort(fmt.Sprintf("%s:%d", s.Addr, s.Port)))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("failed to start udp listener: %s", err)
	}

	go s.listenForMessages(conn)
	log.Printf("Server started. Listening on %d", s.Port)
	s.wg.Add(1)
	s.wg.Wait()
}

func main() {
	server := NewServer("127.0.0.1", ListenPort)

	QuietHours[8] = 0
	QuietHours[9] = 0
	QuietHours[10] = 0
	QuietHours[11] = 0
	QuietHours[12] = 0
	QuietHours[13] = 0
	server.Listen()
}
