package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	connHost = "localhost"
	connPort = "3333"
	connType = "tcp"
)

type Command struct {
	verb  string
	key   string
	value string
}

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
	fmt.Println("Ran")
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		parsedCommand, err := parseCommand(data)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Print(parsedCommand)

	}
	//conn.Close()
}

func parseCommand(rawData string) (Command, error) {
	parsedData := strings.Split(rawData, " ")
	fmt.Printf("len=%d %v\n", len(parsedData), parsedData)

	if len(parsedData) == 2 {
		return Command {
			verb: parsedData[0],
			key: parsedData[1],
		}, nil
	} else if len(parsedData) == 3 {
		return Command {
			verb: parsedData[0],
			key: parsedData[1],
			value: parsedData[2],
		}, nil
	} else {
		return Command{}, errors.New("invalid command format")
	}
}

//Need to parse the string input to understand what to do
//based on the parsed string, run an action
