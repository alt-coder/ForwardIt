package client

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"bytes"
)

func Client( port int, destinationServerAddr string) {
	// Specify the local client port number and the destination server address
	localClientPort := 	port // Replace with your destination server address

	// Create a TCP listener to capture packets from the local client
	listenAddr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for packets on port %d...\n", localClientPort)

	for {
		// Accept incoming client connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go func() {
			// Read the packet content from the connection
			packetData, err := io.ReadAll(conn)
			if err != nil {
				fmt.Println("Error reading packet data:", err)
				conn.Close()
				return
			}

			// Create an HTTP POST request to send the packet data to the server
			client := &http.Client{}
			req, err := http.NewRequest("POST", destinationServerAddr, bytes.NewReader(packetData))
			if err != nil {
				fmt.Println("Error creating HTTP request:", err)
				conn.Close()
				return
			}

			// Add a custom header with the local client port number
			req.Header.Add("X-Local-Client-Port", strconv.Itoa(localClientPort))

			// Send the HTTP request to the destination server
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request to server:", err)
				conn.Close()
				return
			}
			defer resp.Body.Close()

			// Copy the server's response back to the local client
			_, err = io.Copy(conn, resp.Body)
			if err != nil {
				fmt.Println("Error copying response to client:", err)
			}

			conn.Close()
		}()
	}
}
