package main

import (
	"fmt"
	"net"
)

type socket struct{
	connection net.Conn
	err error
}

var s = &socket{}

func establishConnection() (conn net.Conn, err error) {
	s.connection, s.err = net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	
	if s.err != nil {
		fmt.Println(s.err)
	}

	return s.connection, s.err
}

func sendData(numberToSend string) {
	_, s.err = s.connection.Write([]byte(numberToSend))
	buffer := make([]byte, 1024)
	mLen, err := s.connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	
	fmt.Println("Received: ", string(buffer[:mLen]))
}