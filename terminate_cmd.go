package main

import (
	"log"

)

// TerminateCommand exported to add to router
type TerminateCommand struct {}

func (t TerminateCommand) execute(s *Server) {
	log.Printf("Terminate sequence being executed")
	s.stop()
}

func (t TerminateCommand) getPath() string {
	return "terminate"
}