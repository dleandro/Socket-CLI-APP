package main

type TerminateCommand struct {}

func (t TerminateCommand) execute(s *Server) {
	s.Stop()
}

func (t TerminateCommand) getPath() string {
	return "terminate"
}