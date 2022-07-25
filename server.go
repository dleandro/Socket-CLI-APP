package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/netutil"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "4000"
	SERVER_TYPE = "tcp"
	MAXIMUM_NUMBER_OF_CONNECTIONS = 5
	NUMBER_OF_CHARS_PER_MESSAGE = 9
	LOG_FILE_NAME = "numbers.log"
)

var env string
var connectionEstablished = make(chan bool)
var numberMessagesReceived = make(chan *NumberMessages)
var appCleanlyShutdown = make(chan bool)

type Server struct {
	listener net.Listener
	quit chan interface{}
	wg sync.WaitGroup
	cmds []Command
	numberMessages *NumberMessages
}

func main() {
	os.Create(LOG_FILE_NAME)
	CreateServer()
}

func scheduleSummary(s *Server) {
	for range time.Tick(10 * time.Second) {
		s.numberMessages.getSummary()
		s.numberMessages.transferSummary()
		s.numberMessages.resetCurrentSummary()
	}
}

func CreateServer() *Server {
	env = os.Getenv("GO_ENV")
	s := &Server{
		quit: make(chan interface{}),
		cmds: []Command{TerminateCommand{}},
		numberMessages: &NumberMessages{},
	}
	l, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		log.Fatal(err)
	}
	l = netutil.LimitListener(l, MAXIMUM_NUMBER_OF_CONNECTIONS)
	s.listener = l

	log.Println("Server Running...")
	log.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	log.Println("Waiting for client...")
	go scheduleSummary(s)
	s.wg.Add(1)
	s.serve()
	return s
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	if env == "TEST" {
		appCleanlyShutdown <- true
	}
	s.wg.Wait()
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <- s.quit:
				return
			default:
				log.Println("error: ", err)
			}
		} else {
			log.Println("client connected")
			if env == "TEST" {
				connectionEstablished <- true
			}
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
	buf := make([]byte, 128)
ReadLoop:
	for {
		select {
		case <-s.quit:
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				if _, ok := err.(*net.OpError); ok {
					continue ReadLoop
				} else if err != io.EOF {
					log.Println("read error", err)
					return
				}
			}
			if n == 0 {
				return
			}
			input := string(buf[:n - 2])  // subtract by two because the commands always have an extra /n at the end
			err = handleCommand(input, s, conn)
			log.Printf("received from %v: %s", conn.RemoteAddr(), input)
			if (env == "TEST") {
				numberMessagesReceived <- s.numberMessages
			}
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}