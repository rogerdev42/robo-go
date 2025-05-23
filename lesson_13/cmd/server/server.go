package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"lesson_13/internal/commands"
	"log"
	"net"
	"strings"
)

func StartServer(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}()

	log.Printf("Server started on %s", address)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
		log.Printf("Connection closed: %s", conn.RemoteAddr())
	}()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		msg := scanner.Text()
		line, derr := base64.StdEncoding.DecodeString(msg)
		if derr != nil {
			log.Printf("Decode error: %v", derr)
		}
		elems := strings.Split(string(line), " ")

		if len(elems) != 3 {
			if _, err := writer.WriteString("Error: Invalid command\n"); err != nil {
				log.Printf("write error: %v", err)
				return
			}
			if err := writer.Flush(); err != nil {
				log.Printf("flush error: %v", err)
				return
			}
			continue
		}
		storage := elems[0]
		command := elems[1]
		param := elems[2]

		var resp string
		var err error

		switch storage {
		case commands.StoreStorageName:
			resp, err = commands.ExecStore(command, param)
		case commands.ColStorageName:
			resp, err = commands.ExecCol(command, param)
		default:
			err = fmt.Errorf("unknown command: %s", storage)
		}

		if err != nil {

			r := struct {
				Status string `json:"status"`
				Value  string `json:"value"`
			}{
				Status: "error",
				Value:  err.Error(),
			}

			resp, merr := json.Marshal(r)
			if merr != nil {
				log.Println("internal error:", merr)
				return
			}

			if _, werr := writer.WriteString(string(resp) + "\n"); werr != nil {
				log.Printf("write error: %v", werr)
				return
			}
			if ferr := writer.Flush(); ferr != nil {
				log.Printf("flush error: %v", ferr)
				return
			}
			continue
		}

		if _, werr := writer.WriteString(resp + "\n"); werr != nil {
			log.Printf("Write error: %v", werr)
			return
		}
		if ferr := writer.Flush(); ferr != nil {
			log.Printf("flush error: %v", ferr)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v", err)
	}
}
