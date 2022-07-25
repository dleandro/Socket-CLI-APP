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
	
	socket := establishConnection()
	defer socket.connection.Close()
	
	<- connectionEstablished
	
	t.Run("send message through socket_message is received", func(f *testing.T) {
		numberToSend := "123456789"
		
		socket.sendData(numberToSend)
		
		numberMessages := <- numberMessagesReceived
		assert.Equal(t, numberToSend, numberMessages.messageList[0])
	})
	
	t.Run("send message with more than 9 digits_message isn't saved to queue", func (t *testing.T)  {
		numberToSend := "1234567892222"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		
		numberMessages := <- numberMessagesReceived

		for _, msg := range numberMessages.messageList {
			if msg == numberToSend {
				assert.Fail(t, "number was processed")
			} 
		}
		assert.True(t, true)
	})
	
	t.Run("send message with less than 9 digits_message isn't saved to queue", func(t *testing.T) {
		numberToSend := "1"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		numberMessages := <- numberMessagesReceived

		for _, msg := range numberMessages.messageList {
			if msg == numberToSend {
				assert.Fail(t, "number was processed")
			} 
		}
		assert.True(t, true)
	})
	
	t.Run("send message with letters and numbers_message isn't saved to queue", func(t *testing.T) {
		numberToSend := "3*567CEW2"
		socket := establishConnection()
		defer socket.connection.Close()
		
		<- connectionEstablished
		socket.sendData(numberToSend)
		
		numberMessages := <- numberMessagesReceived

		for _, msg := range numberMessages.messageList {
			if msg == numberToSend {
				assert.Fail(t, "number was processed")
			} 
		}
		assert.True(t, true)
	})

	t.Run("establish 6 connections and send message_6th connection isn't established", func(f *testing.T) {
		var socket *Socket
		numberToSend := "987654321"
		for i := 0; i < 6; i++ {
			socket = establishConnection()
		}
		defer socket.connection.Close()
		
		socket.sendData(numberToSend)
		
		select {
		case numberMessages := <- numberMessagesReceived:
			for _, msg := range numberMessages.messageList {
				if msg == numberToSend {
					assert.Fail(t, "number was processed")
				} 
			}
			assert.True(t, true)
		case <- time.After(3 * time.Second):
			assert.True(t, true)
		}		
		
	})
}

func TestTerminateSequence(t *testing.T) {
	os.Setenv("GO_ENV", "TEST")
	terminateStr := "terminate"
	
	go CreateServer()
	
	socket := establishConnection()
	defer socket.connection.Close()
	
	<- connectionEstablished
	socket.sendData(terminateStr)

	select {
	case <- appCleanlyShutdown:
		assert.True(t, true)
	case <- time.After(10 * time.Second):
		assert.Fail(t, "app wasn't cleanly shutdown")
	}
}


func TestAppCapacity(t * testing.T) {
	var numberMessages *NumberMessages

	os.Setenv("GO_ENV", "TEST")
	
	go CreateServer()

	socket := establishConnection()
	defer socket.connection.Close()
	

	<- connectionEstablished
	for i := 0; i < 400000; i++ {
		socket.sendData("123456789")
		numberMessages = <- numberMessagesReceived
	}

	assert.GreaterOrEqual(t, len(numberMessages.messageList), 400000)
	
}