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

var dataStore = make(map[string]string)

func main() {
	//Listen for connections
	dStream, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println(err)
		//TODO send error information back to client
		return
	}
	fmt.Println("Database server running on " + connHost + ":" + connPort)
	defer dStream.Close()

	//Accepts connection requests
	for {
		conn, err := dStream.Accept()
		fmt.Println("User connected")
		if err != nil {
			fmt.Println(err)
			//TODO send error information back to client
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
			//TODO send error information back to client
			return
		}

		parsedCommand, err := parseCommand(data)
		if err != nil {
			fmt.Println(err)
			//TODO send error information back to client
		}
		fmt.Printf("verb: %s, key: %s, value: %s\n",parsedCommand.verb, parsedCommand.key, parsedCommand.value)

		handleCommand(parsedCommand)
		fmt.Println(dataStore)

	}
	//conn.Close()
}

func parseCommand(rawData string) (Command, error) {
	parsedData := strings.Split(rawData, " ")

	if len(parsedData) == 2 {
		return Command{
			verb: parsedData[0],
			key:  parsedData[1],
			value: "",
		}, nil
	} else if len(parsedData) > 2 {

		tmpSlice := []string{}
		for x := 2; x < len(parsedData); x++ {
			tmpSlice = append(tmpSlice, parsedData[x])
		}
		tmpStr := strings.Join(tmpSlice, " ")

		return Command{
			verb:  parsedData[0],
			key:   parsedData[1],
			value: tmpStr,
		}, nil
	} else {
		return Command{}, errors.New("invalid command format")
	}
}

func handleCommand(cmd Command) error {
	//GET, SET, DELETE, UPDATE
	verb := strings.ToUpper(cmd.verb)

	switch verb {
	case "GET":
		getCommand()
		return nil
	case "SET":
		setCommand(cmd)
		return nil
	case "DELETE":
		deleteCommand()
		return nil
	case "UPDATE":
		updateCommand()
		return nil
	default:
		return errors.New("invalid command, use GET, SET, UPDATE or DELETE")
	}
}

func getCommand() {
	
}

func setCommand(cmd Command) {
	dataStore[cmd.key] = cmd.value
	
}

func deleteCommand() {

}

func updateCommand() {

}

//Need to parse the string input to understand what to do
//based on the parsed string, run an action
