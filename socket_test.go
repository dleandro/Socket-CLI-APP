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
	
	t.Run("client sends two valid messages_both get saved", func(t *testing.T) {
		m1 := "123456789"
		m2 := "987654321"
		socket := establishConnection()
		defer socket.connection.Close()
		
		
		<- connectionEstablished
		socket.sendData(m1)
		socket.sendData(m2)
		
		numberMessages := <- numberMessagesReceived
		assert.Equal(t, m1, numberMessages.messageList[0])
		assert.Equal(t, m2, numberMessages.messageList[1])
	})
	
	t.Run("establish 6 connections and send message_6th connection message isn't received", func(f *testing.T) {
		numberToSend := "123242444"
		var socket *Socket
		for i := 0; i < 6; i++ {
			socket = establishConnection()
		}
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		select {
		case <- connectionEstablished:
			assert.True(t, true)
		case <- time.After(3 * time.Second):
			assert.Fail(t, "connection wasn't established in time")
		}
		
	})
	
	t.Run("send message with more than 9 digits_message isn't saved to queue", func (t *testing.T)  {
		numberToSend := "1234567892222"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		numberMessages := <- numberMessagesReceived
		assert.NotEqual(t, numberToSend, numberMessages.messageList[0])
	})
	
	t.Run("send message with less than 9 digits_message isn't saved to queue", func(t *testing.T) {
		numberToSend := "1"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		numberMessages := <- numberMessagesReceived
		assert.NotEqual(t, numberToSend, numberMessages.messageList[0])
	})
	
	t.Run("send message with letters and numbers_message isn't saved to queue", func(t *testing.T) {
		numberToSend := "3*567CEW2"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		numberMessages := <- numberMessagesReceived
		assert.NotEqual(t, numberToSend, numberMessages.messageList[0])
	})
	
	t.Run("client sends a message with the word terminate followed by a server-native newline sequence_all clients are disconnected and app shuts down", func(t *testing.T) {
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData("terminate")

		select {
		case <- appCleanlyShutdown:
			assert.True(t, true)
		case <- time.After(3 * time.Second):
			assert.Fail(t, "app was cleanly shutdown")
		}
	})
}




