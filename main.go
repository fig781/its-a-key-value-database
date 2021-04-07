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

		normalizedData := normaliseData(data)

		parsedCommand, err := parseCommand(normalizedData)
		if err != nil {
			fmt.Println(err)
			sendResponse(conn, parsedCommand.key, err)
		} else {
			fmt.Printf("Verb: %s, Key: %s, Value: %s\n", parsedCommand.verb, parsedCommand.key, parsedCommand.value)
			//Every input past this should be valid, if an input is not valid it should throw an error.

			value, err := handleCommand(parsedCommand)
			if err != nil {
				fmt.Println(err)
			}
			sendResponse(conn, value, err)

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

func normaliseData(data string) string {
	var runeSlice []rune
	for _, c := range data {
		//ignors all ASCII from 0-31
		if c <= 31 || c == 127 {
			//Remove backspace char
			if len(runeSlice) != 0 && c == 8 {
				runeSlice = runeSlice[:len(runeSlice)-1]
			}
		} else {
			runeSlice = append(runeSlice, c)
		}
	}

	returnStr := string(runeSlice)
	returnStr = strings.TrimSuffix(returnStr, "\r\n")
	return returnStr
}

func parseCommand(rawData string) (Command, error) {
	parsedData := strings.Split(rawData, " ")

	if len(parsedData) == 1 {
		return Command{
			verb:  parsedData[0],
			key:   "",
			value: "",
		}, nil
	} else if len(parsedData) == 2 {
		return Command{
			verb:  parsedData[0],
			key:   parsedData[1],
			value: "",
		}, nil
	} else {

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
	case "GETALL":
		value, err = getAllCommand(cmd)
	case "LEN":
		value, err = lenCommand(cmd)
	case "GETKEYS":
		value, err = getKeysCommand(cmd)
	case "GETVALUES":
		value, err = getValuesCommand(cmd)
	case "EXISTS":
		value, err = existsCommand(cmd)
	case "PING":
		value = pingCommand(cmd)
	default:
		value = "-ERR"
		err = errors.New("invalid command")
	}

	return value, err
}

func getCommand(cmd Command) (string, error) {
	val, exists := dataStore[cmd.key]
	if exists {
		//"+Aden\r\n"
		formatedReturnVal := "+" + val + "\r\n"
		return formatedReturnVal, nil
	} else {
		return "-ERR", errors.New("key does not exists\r\n")
	}
}

func setCommand(cmd Command) (string, error) {
	_, exists := dataStore[cmd.key]
	if !exists {
		dataStore[cmd.key] = cmd.value
		return "+OK\r\n", nil
	} else {
		return "-ERR", errors.New("key already exists\r\n")
	}
}

func updateCommand(cmd Command) (string, error) {
	_, exists := dataStore[cmd.key]
	if exists {
		dataStore[cmd.key] = cmd.value
		return "+OK\r\n", nil
	} else {
		return "-ERR", errors.New("key does not exist\r\n")
	}
}

func deleteCommand(cmd Command) (string, error) {
	_, exists := dataStore[cmd.key]
	if exists {
		delete(dataStore, cmd.key)
		return "+OK\r\n", nil
	} else {
		return "-ERR", errors.New("key does not exist\r\n")
	}
}

func getAllCommand(cmd Command) (string, error) {
	if len(dataStore) != 0 {
		allEntries := "+"
		// "+user1\r\nAden\r\nuser2\r\nEilers\r\n"
		for key, value := range dataStore {
			allEntries += key + "\r\n" + value + "\r\n"
		}
		return allEntries, nil
	} else {
		return "-ERR", errors.New("no entries in database\r\n")
	}
}

func lenCommand(cmd Command) (string, error) {
	return "+" + fmt.Sprint(len(dataStore)) + "\r\n", nil
}

func getValuesCommand(cmd Command) (string, error) {
	if len(dataStore) != 0 {
		allValues := "+"
		// "+Aden\r\nEilers\r\n"
		for _, value := range dataStore {
			allValues += value + "\r\n"
		}
		return allValues, nil
	} else {
		return "-ERR", errors.New("no entries in database\r\n")
	}
}

func getKeysCommand(cmd Command) (string, error) {
	if len(dataStore) != 0 {
		allKeys := "+"
		// "+user1\r\nuser2\r\n"
		for key := range dataStore {
			allKeys += key + "\r\n"
		}
		return allKeys, nil
	} else {
		return "-ERR", errors.New("no entries in database\r\n")
	}
}

func existsCommand(cmd Command) (string, error) {
	_, exists := dataStore[cmd.key]
	if exists {
		return "+1\r\n", nil
	} else {
		return "+0\r\n", nil
	}
}

func pingCommand(cmd Command) string {
	return "+PONG"
}
