package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	address = "127.0.0.1"
	port    = 6379
)

func main() {
	send, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	defer send.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		var content string
		fmt.Print("input message: ")
		_, _ = fmt.Scan(&content)
		contentList := []byte(content)
		bodyLength := len(contentList)
		data := make([]byte, 4)
		binary.BigEndian.PutUint32(data[0:4], uint32(bodyLength))
		data = append(data, contentList...)
		fmt.Println(data)
		_, err = send.Write(data)
		if err != nil {
			fmt.Println(err)
		}
		response := make([]byte, 100)
		send.Read(response)
		fmt.Println(string(response))
	}
}
