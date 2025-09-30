package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

// Reads lines from a stream and emits them into a channel
func getLinesFromReader(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8) // read in 8-byte chunks
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

		if len(str) != 0 {
			out <- str
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go func(c net.Conn) {
			defer c.Close()
			for line := range getLinesFromReader(c) {
				fmt.Printf("read: %s\n", line)
			}
		}(conn)
	}
}
