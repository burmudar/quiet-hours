package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/burmudar/quiet-hours/models"
)

const (
	ListenPort  = 20111
	CommandPort = 20112
)

var QuietHours = []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

type Server struct {
	closer io.Closer
	Addr   string
	Port   int
	wg     sync.WaitGroup
	done   chan any
}

func New(addr string, port int) *Server {
	return &Server{
		Addr: addr,
		Port: port,
		wg:   sync.WaitGroup{},
	}
}

func (s *Server) Shutdown() { close(s.done) }

func handleQuery(p models.Preamble, data []byte) ([]byte, error) {
	var q models.QuietQuery
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

	resp := models.QuietResponse{
		WakeUpHour:  uint(wakeUpIn),
		IsQuietTime: inQuietHours,
		Whoru:       q.Whoami,
	}

	log.Printf("%+v", resp)

	return resp.Marshal()
}

func (s *Server) listenForMessages(conn net.PacketConn) {
	process := func(conn net.PacketConn, packet []byte, addr net.Addr) {
		log.Printf("[INFO] processing packet with size %d", len(packet))
		log.Printf("[DEBUG] received: %v", packet)
		p := &models.Preamble{}
		n, err := p.Unmarshal(packet)
		if err != nil {
			log.Printf("[ERR] error processing preamble: %s", err)
			return
		}

		switch p.PacketType {
		case models.QuietQueryType:
			if data, err := handleQuery(*p, packet[n:]); err != nil {
				log.Printf("[ERR] error processing query: %s", err)
			} else {
				log.Printf("[INFO] resp: %v", data)
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
					log.Printf("[INFO] Read %d from Conn", n)
					go process(conn, packet[:n], addr)
				} else if err != nil {
					log.Printf("[ERR] error processing packet: %s", err)
					return
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
	defer conn.Close()

	go s.listenForMessages(conn)
	log.Printf("Server started. Listening on %d", s.Port)
	s.wg.Add(1)
	s.wg.Wait()
}
