package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesFromReader(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

	str := ""
	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			break

	}()
	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for connection, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
	}

	for line := range getLinesFromReader(connection) {
		fmt.Printf("read: %s\n", line)
	}

	str := ""
	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			break
		}

		data = data[:n]
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			str += string(data[:i])
			data = data[i+1:]
			out <- str
			str = ""
		}
		str += string(data)
	}

	if leng(str) != 0 {
		fmt.Printf("read: %s\n", str)
	}
}
