package main

import (
	"errors"
	"net"
	"unicode"

	m "tcp-socket-app/messages"
)

// Command interface used to establish cmds in the router
type Command interface {
	getPath() string
	execute(s *Server)
}

func handleCommand(input string, s *Server, c net.Conn) error {
	for _, cmd := range s.cmds {
		if cmd.getPath() == input {
			cmd.execute(s)
			return nil
		}
	}

	return s.numberMessages.HandleMessage(input)
}
