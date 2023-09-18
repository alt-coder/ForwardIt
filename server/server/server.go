package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

func Server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the custom header X-Local-Client-Port
		localClientPortHeader := r.Header.Get("X-Local-Client-Port")
		localClientPort, err := strconv.Atoi(localClientPortHeader)
		if err != nil {
			http.Error(w, "Invalid X-Local-Client-Port header", http.StatusBadRequest)
			return
		}

		// Read the request body
		bodyData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		// Forward the request body to the local port
		err = forwardToPort(localClientPort, bodyData)
		if err != nil {
			http.Error(w, "Error forwarding request to local port", http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Request forwarded to local port %d", localClientPort)
	})

	// Start the HTTP server
	fmt.Println("Server listening on :8080...")
	http.ListenAndServe(":8080", nil)
}

func forwardToPort(port int, data []byte) error {
	// Create a TCP connection to the local port
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the data to the local port
	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
