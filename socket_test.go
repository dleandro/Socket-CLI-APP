package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSocket(t *testing.T) {
	
	os.Setenv("GO_ENV", "TEST")
	
	go CreateServer()
	
	t.Run("send message through socket_message is received", func(f *testing.T) {
		numberToSend := "123456789"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		numberMessages := <- numberMessagesReceived
		assert.Equal(t, numberToSend, numberMessages.messageList[0])
	})
	
	t.Run("establish 6 connections and send message_6th connection isn't established", func(f *testing.T) {
		var socket *Socket
		for i := 0; i < 6; i++ {
			socket = establishConnection()
			defer socket.connection.Close()
		}
		
		<- connectionEstablished
		socket.sendData("123456789")
		
		select {
		case <- numberMessagesReceived:
			assert.Fail(t, "message was received")
		case <- time.After(3 * time.Second):
			assert.True(t, true)
		}
		
	})
	
	t.Run("send message with more than 9 digits_message isn't saved to queue", func (t *testing.T)  {
		numberToSend := "1234567892222"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		select {
		case <- numberMessagesReceived:
			assert.Fail(t, "number was processed")
		case <- time.After(3 * time.Second):
			assert.True(t, true)
		}
	})
	
	t.Run("send message with less than 9 digits_message isn't saved to queue", func(t *testing.T) {
		numberToSend := "1"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		select {
		case <- numberMessagesReceived:
			assert.Fail(t, "number was processed")
		case <- time.After(3 * time.Second):
			assert.True(t, true)
		}
	})
	
	t.Run("send message with letters and numbers_message isn't saved to queue", func(t *testing.T) {
		numberToSend := "3*567CEW2"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		select {
		case <- numberMessagesReceived:
			assert.Fail(t, "number was processed")
		case <- time.After(3 * time.Second):
			assert.True(t, true)
		}
	})
	
	t.Run("client sends a message with the word terminate followed by a server-native newline sequence_all clients are disconnected and app shuts down", func(t *testing.T) {
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData("terminate")

		select {
		case <- appCleanlyShutdown:
			assert.True(t, true)
		case <- time.After(10 * time.Second):
			assert.Fail(t, "app wasn't cleanly shutdown")
		}
	})

	t.Run("client sends 2 million messages in a 10 second period-messages are processed", func(t *testing.T) {
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		for i := 0; i < 200000; i++ {
			socket.sendData("123456789")
		}

		// break it into smaller cycles with go routines and probably use random numbers idk

		numberMessages := <- numberMessagesReceived

		time.AfterFunc(10 * time.Second, func() {})

		assert.Equal(t, 2000, len(numberMessages.messageList))
	})
	
}




