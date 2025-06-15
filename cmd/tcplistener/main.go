package main

import (
	"errors" // Provides functions for error handling
	"fmt"    // Implements formatted I/O functions
	"io"     // Provides basic I/O primitives
	"log"    // Provides logging functions

	// Provides functions for OS-level operations
	"net"
	"strings"
)

// What port to lisen on
const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		// Holds the current line
		currentLineContents := ""
		// Infinite loop to read the file in chunks
		for {
			// Create a buffer to hold the data read from the file
			buffer := make([]byte, 8, 8)
			// Read data into the buffer
			n, err := f.Read(buffer)
			if err != nil {
				// This sees if the current line is empty
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				// Check if the end of the file has been reached
				if errors.Is(err, io.EOF) {
					break // Exit the loop if the end of the file is reached
				}
				// If another error occurs, print the error and exit the loop
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			// Convert the read bytes into a string and print it
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
