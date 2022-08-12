package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/net/netutil"
	logs "tcp-socket-app/logs"
	m "tcp-socket-app/messages"
)

const (
	serverHost                 = "localhost"
	serverPort                 = "4000"
	serverType                 = "tcp"
	maximumNumberOfConnections = 5
)

var env string
var connectionEstablished = make(chan bool)
var appCleanlyShutdown = make(chan bool)

// Server obj
type Server struct {
	listener       net.Listener
	quit           chan interface{}
	wg             sync.WaitGroup
	cmds           []Command
	numberMessages *m.NumberMessages
}

func main() {
	logger := logs.InitLogger()
	CreateServer(logger.LogFileIsSetupChan)
}

// CreateServer Exported function to initialize server
func CreateServer(LogFileIsSetupChan chan (bool)) *Server {
	env = os.Getenv("GO_ENV")
	s := &Server{
		quit: make(chan interface{}),
		cmds: []Command{TerminateCommand{}},
	}
	s.numberMessages = m.Init(s.quit, s.wg)
	l, err := net.Listen(serverType, serverHost+":"+serverPort)
	if err != nil {
		log.Fatal(err)
	}
	l = netutil.LimitListener(l, maximumNumberOfConnections)
	s.listener = l

	log.Println("Server Running...")
	log.Println("Listening on " + serverHost + ":" + serverPort)
	log.Println("Waiting for client...")

	<-LogFileIsSetupChan
	s.wg.Add(1)
	s.serve()
	return s
}

func (s *Server) stop() {
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
			case <-s.quit:
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

		// TODO: case for receiving a message from other client connections could add a new channel
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
			input := string(buf[:n-2]) // subtract by two because the commands always have an extra /n at the end
			log.Printf("received from %v: %s", conn.RemoteAddr(), input)
			err = handleCommand(input, s, conn) // TODO: trigger channel for chat
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
