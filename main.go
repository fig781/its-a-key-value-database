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
		connPort = "3332"
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
			return
		}
		fmt.Println("Client connected")

		go handleConnection(connection)
	}
}

func handleConnection(conn net.Conn) {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			sendResponse(conn, "", err)
			return
		}
		data = strings.TrimSuffix(data, "\r\n")

		parsedCommand, err := parseCommand(data)
		if err != nil {
			fmt.Println(err)
			sendResponse(conn, parsedCommand.key, err)
		} else {
			fmt.Printf("Verb: %s, Key: %s, Value: %s\n", parsedCommand.verb, parsedCommand.key, parsedCommand.value)
			//Every input past this should be valid, if an input is not valid it should throw an error.

			value, err := handleCommand(parsedCommand)
			if err != nil {
				fmt.Println(err)
				sendResponse(conn, value, err)
			} else {
				sendResponse(conn, value, err)
			}
			fmt.Println("Datastore:", dataStore)
		}
	}
}

func sendResponse(conn net.Conn, res string, err error) {
	conn.Write([]byte(constructResponse(res, err)))
}

func constructResponse(value string, err error) string {
	if err == nil {
		return value
	} else {
		return value + " " + err.Error()
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
		return Command{
			verb:  "ERR",
			key:   "",
			value: "",
		}, errors.New("invalid command format, needs to be <verb> <key> <value>")
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
		err = errors.New("invalid command, use GET, SET, UPDATE or DELETE")
	}

	return value, err
}

func getCommand(cmd Command) (string, error) {
	val, exists := dataStore[cmd.key]
	if exists {
		return val, nil
	} else {
		return "ERR", errors.New("key does not exists")
	}
}

func setCommand(cmd Command) (string, error) {
	_, exists := dataStore[cmd.key]
	if !exists {
		dataStore[cmd.key] = cmd.value
		return "OK", nil
	} else {
		return "ERR", errors.New("key already exists")
	}
}

func updateCommand(cmd Command) (string, error) {
	_, exists := dataStore[cmd.key]
	if exists {
		dataStore[cmd.key] = cmd.value
		return "OK", nil
	} else {
		return "ERR", errors.New("key does not exist")
	}
}

func deleteCommand(cmd Command) (string, error) {
	_, exists := dataStore[cmd.key]
	if exists {
		delete(dataStore, cmd.key)
		return "OK", nil
	} else {
		return "ERR", errors.New("key does not exist")
	}
}
