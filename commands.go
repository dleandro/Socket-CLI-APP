package main

import (
	"net"
	"errors"
	"unicode"
)

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
	
	err := checkIfMessageIsValid(input)
	
	if err != nil {
		return err
	}
	
	return s.numberMessages.handleNewNumber(input)
}

func checkIfMessageIsValid (input string) error {
	// check for numberOfDigits
	if len(input) != NUMBER_OF_CHARS_PER_MESSAGE {
		return errors.New("input message has incorrect number of digits")
	}
	
	// check if every char is a digit
	for _, char := range input {
		if !unicode.IsDigit(char) {
			return errors.New("input message has invalid chars")
		}
	}
	
	return nil
}
