package main

import (
	"log"
)

type TerminateCommand struct {}

func (t TerminateCommand) execute(s *Server) {
	log.Printf("Terminate sequence being executed")
	s.Stop()
}

func (t TerminateCommand) getPath() string {
	return "terminate"
}