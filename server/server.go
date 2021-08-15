package server

import (
	"io"
	"math"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	quit     chan interface{}
	listener net.Listener
	wg       sync.WaitGroup
	count    int64
	mapper   map[string]int64
	mu       sync.RWMutex
	logger   *log.Entry
}

func getPattern() map[string]int64 {
	return map[string]int64{"q": 0, "w": 0, "e": 0, "r": 0, "t": 0, "y": 0}
}

func NewServer(addr string, logger *log.Entry) *Server {
	s := &Server{
		quit:   make(chan interface{}),
		mapper: getPattern(),
		mu:     sync.RWMutex{},
		logger: logger,
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("listening on: ", addr)
	s.listener = l
	s.wg.Add(1)

	go s.serve()
	return s
}

func (s *Server) Debug() {
	s.logger.Infof("current counter is: %d", s.count)
}

func (s *Server) Stop() error {
	close(s.quit)
	err := s.listener.Close()
	s.wg.Wait()
	return err
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				s.logger.Error("accept error: ", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleConection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1)
	mapper := getPattern()

	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			s.logger.Println("read error", err)
			return
		}
		if n == 0 {
			break
		}
		key := string(buf)
		if val, ok := mapper[key]; ok {
			mapper[key] = val + 1
		}
	}

	s.mu.Lock()
	count := int64(math.MaxInt64)
	for char, val := range mapper {
		prev := s.mapper[char]
		if val+prev < count {
			count = val + prev
		}
	}
	s.mapper = mapper
	s.count = s.count + count
	s.mu.Unlock()
}
