package main

import (
        "fmt"
        "net"
        "os"
)

const (
        SERVER_HOST = "localhost"
        SERVER_PORT = "4000"
        SERVER_TYPE = "tcp"
		MAXIMUM_NUMBER_OF_CONNECTIONS = 5
)


var messageList []string
var currentConnections = 0


func main() {
        fmt.Println("Server Running...")
        server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
        if err != nil {
                fmt.Println("Error listening:", err.Error())
                os.Exit(1)
        }
        defer server.Close()
        fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
        fmt.Println("Waiting for client...")
        for currentConnections <= MAXIMUM_NUMBER_OF_CONNECTIONS {
                connection, err := server.Accept()
                if err != nil {
                        fmt.Println("Error accepting: ", err.Error())
                        os.Exit(1)
                }
				currentConnections++
                fmt.Println("client connected")
                go processClient(connection)
        }
}

func processClient(connection net.Conn) {
        buffer := make([]byte, 1024)
        mLen, err := connection.Read(buffer)
        if err != nil {
                fmt.Println("Error reading:", err.Error())
        }
        fmt.Println("Received: ", string(buffer[:mLen]))
		messageList = append(messageList, string(buffer[:mLen]))

        _, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
        connection.Close()
}