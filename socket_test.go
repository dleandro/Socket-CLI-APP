package main

import (
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSocket(t *testing.T) {

	go main()

	establishConnection()
	
	t.Run("send message through socket_message is received", func(f *testing.T) {
		numberToSend := "123456789"
		sendData(numberToSend)

		assert.Equal(t, numberToSend, messageList[len(messageList)-1])

	})

	t.Run("establish 6 connections_6th connection fails", func(f *testing.T) {
		var c net.Conn
		for i := 0; i < 5; i++ {
			c, _ = establishConnection()
			log.Println(currentConnections)
		}

		assert.Equal(t, nil, c)

	})

}