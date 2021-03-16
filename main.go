package main

import (
	"bufio"
	"fmt"
	"net"
)

const (
	connHost = "localhost"
	connPort = "3333"
	connType = "tcp"
)

func main() {
	//Listen for connections
	dStream, err := net.Listen(connType, connHost + ":" + connPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Database server running on " + connHost + ":" + connPort)
	defer dStream.Close()

	//Accepts connection requests
	for {
		conn, err := dStream.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(conn)
	}
}

//Handles data that streams in from TCP server connection
func handleConnection(conn net.Conn) {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(data)
	}
	//conn.Close()
}