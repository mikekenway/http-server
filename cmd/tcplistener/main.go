package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// Convert a readable stream into a channel of strings,
// where each string represents one line of input.
func getDataStream(reader io.ReadCloser) <-chan string {
	dataChannel := make(chan string, 1)

	go func() {
		defer reader.Close()
		defer close(dataChannel)

		var dataBuffer strings.Builder

		for {
			// Pull the next chunk of bytes from the stream
			data := make([]byte, 8)
			dataRead, err := reader.Read(data)

			// Stop reading if the stream closes or fails
			if err != nil {
				break
			}

			// Focus only on the actual bytes read this time
			data = data[:dataRead]

			// If this chunk includes a new line, that marks the end of a message
			if dataStream := bytes.IndexByte(data, '\n'); dataStream != -1 {
				dataBuffer.Write(data[:dataStream]) // complete the line
				dataChannel <- dataBuffer.String()  // hand it off
				dataBuffer.Reset()                  // prepare for the next
				data = data[dataStream+1:]          // keep remainder after newline
			}

			// Add whatever is left to the current message-in-progress
			if len(data) > 0 {
				dataBuffer.Write(data)
			}
		}

		// If the stream ends mid-line, send the unfinished line too
		if dataBuffer.Len() != 0 {
			dataChannel <- dataBuffer.String()
		}
	}()

	return dataChannel
}

func main() {
	// Start a TCP server listening on port 42069
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("listener error:", err)
	}
	defer listener.Close()

	for {
		// Accept a new client connection
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		// Spin up a goroutine to handle this client separately
		go func(connection net.Conn) {
			defer connection.Close()

			// Print each line received from the client
			for line := range getDataStream(connection) {
				fmt.Printf("read: %s\n", line)
			}
		}(connection)
	}
}
