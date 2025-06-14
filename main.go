package main

import (
	"errors" // Provides functions for error handling
	"fmt"    // Implements formatted I/O functions
	"io"     // Provides basic I/O primitives
	"log"    // Provides logging functions
	"os"     // Provides functions for OS-level operations
	"strings"
)

// Define the path to the input file
const inputFilePath = "messages.txt"

func main() {
	// Attempt to open the specified file for reading
	f, err := os.Open(inputFilePath)
	if err != nil {
		// If an error occurs, log the error and terminate the program
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	// Print a message indicating that data is being read from the file
	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	lineschan := getLinesChannel(f)

	for line := range lineschan {
		fmt.Println("read:", line)
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
