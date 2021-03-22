package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type Command struct {
	verb  string
	key   string
	value string
}

var dataStore = make(map[string]string)

func main() {
	const (
		connHost = "localhost"
		connPort = "3333"
		connType = "tcp"
	)

	server, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	fmt.Println("Database server running on " + connHost + ":" + connPort)
	defer server.Close()

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err)
			//TODO send error information back to client
			return
		}
		fmt.Println("Client connected")

		go handleConnection(connection)
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
		fmt.Printf("verb: %s, key: %s, value: %s\n", parsedCommand.verb, parsedCommand.key, parsedCommand.value)
		//Every input past this should be valid, if an input is not valid it should throw an error.

		value, err := handleCommand(parsedCommand)
		if err != nil {
			fmt.Println(err)
			//TODO send error information to the client
		}
		fmt.Println("Value:", value)
		fmt.Println("Datastore:", dataStore)

	}
}

func parseCommand(rawData string) (Command, error) {
	parsedData := strings.Split(rawData, " ")

	if len(parsedData) == 2 {
		return Command{
			verb:  parsedData[0],
			key:   parsedData[1],
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

func handleCommand(cmd Command) (value string, err error) {
	value = ""
	err = nil
	verb := strings.ToUpper(cmd.verb)

	switch verb {
	case "GET":
		value, err = getCommand(cmd)
	case "SET":
		value, err = setCommand(cmd)
	case "UPDATE":
		value, err = updateCommand(cmd)
	case "DELETE":
		value, err = deleteCommand(cmd)
	default:
		err = errors.New("invalid command, use GET, SET or DELETE")
	}

	return value, err
}

func getCommand(cmd Command) (string, error) {
	//cmd.key is problem works fine when it is a hard coded string value
	val, exists := dataStore[cmd.key]
	fmt.Printf("cmd.key type: %T \n", cmd.key)
	fmt.Printf("Val:%s, Exists:%t\n", val, exists)
	if exists {
		return val, nil
	} else {
		return val, errors.New("key does not exist in database")
	}
}

func setCommand(cmd Command) (string, error) {
	dataStore[cmd.key] = cmd.value
	return "OK", nil
}

func updateCommand(cmd Command) (string, error) {
	return "OK", nil
}

func deleteCommand(cmd Command) (string, error) {
	delete(dataStore, cmd.key)
	return "OK", nil
}
