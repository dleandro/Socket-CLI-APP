package main

import (
	"log"
	"net"
)

// Socket object exported for testing reasons
type Socket struct{
	connection net.Conn
	err error
}

var newS = &Socket{}

func establishConnection() *Socket {
	newS.connection, newS.err = net.Dial(serverType, serverHost+":"+serverPort)
	
	if newS.err != nil {
		log.Println(newS.err)
	}

	return newS
}

func (s *Socket) sendData(numberToSend string) {
	_, s.err = s.connection.Write([]byte(numberToSend + "/n")) // add /n to simulate a command entered in the terminal
	if s.err != nil {
		log.Println("Error writing:", s.err.Error())
		return
	}
}

func (s *Socket) send2MMsgs(numberToSend string) {
	for i := 0; i < 200000; i++ {
		_, s.err = s.connection.Write([]byte(numberToSend + "/n")) // add /n to simulate a command entered in the terminal
		if s.err != nil {
			log.Println("Error writing:", s.err.Error())
			return
		}
	}
}