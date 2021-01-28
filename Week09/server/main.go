package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	address = "127.0.0.1"
	port    = 6379
)

func readMessageFunc(conn net.Conn, ch chan<- []byte) {
	reader := bufio.NewReader(conn)
	fmt.Println(conn.RemoteAddr())
	for {
		header, err := reader.Peek(4)
		if err != nil {
			close(ch)
			conn.Close()
			fmt.Println("Server close connection ")
			break
		}
		headerSize := binary.BigEndian.Uint32(header)
		pkg, err := reader.Peek(int(headerSize) + 4)
		if err != nil {
			close(ch)
			conn.Close()
			fmt.Println("Server close connection ")
			break
		}
		body := pkg[4:]
		ch <- body
		reader.Reset(conn)
	}
}

func writeMessageFunc(conn net.Conn, ch <-chan []byte) {
	writer := bufio.NewWriter(conn)
	defer conn.Close()
	for msg := range ch {
		head := make([]byte, 4)
		response := append([]byte("response: "), msg...)
		binary.BigEndian.PutUint32(head[0:4], uint32(len(response)))
		writer.Write(head)
		writer.Write(response)
		writer.Flush()
		fmt.Println("success")
	}
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		ch := make(chan []byte, 10)
		go readMessageFunc(conn, ch)
		go writeMessageFunc(conn, ch)
	}
}
