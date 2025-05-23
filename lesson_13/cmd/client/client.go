package main

import (
	"bufio"
	"encoding/base64"
	"log"
	"net"
	"os"
)

func StartClient(address string) {
	us := bufio.NewScanner(os.Stdin)
	uw := bufio.NewWriter(os.Stdout)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("failed to connect to server:", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("failed to close connection:", err)
		}
	}()

	cw := bufio.NewWriter(conn)
	cr := bufio.NewReader(conn)

	for us.Scan() {
		line := us.Text()
		msg := base64.StdEncoding.EncodeToString([]byte(line))
		_, err := cw.WriteString(msg + "\n")
		if err != nil {
			log.Println("failed to send data to server:", err)
			break
		}
		if err := cw.Flush(); err != nil {
			log.Println("failed to flush data to server:", err)
			break
		}
		resp, err := cr.ReadString('\n')
		if err != nil {
			log.Println("failed to read response from server:", err)
			break
		}
		_, err = uw.WriteString(resp)
		if err != nil {
			log.Println("failed to write response to console:", err)
			break
		}
		if err := uw.Flush(); err != nil {
			log.Println("failed to flush response to console:", err)
			break
		}
	}
	if err := us.Err(); err != nil {
		log.Println("failed to read from console:", err)
	}
}
