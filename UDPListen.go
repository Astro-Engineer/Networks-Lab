package main

import (
	"fmt"
	"net"
)

func main() {
	listenIP := "0.0.0.0" // Listen on all available network interfaces
	listenPort := 12345

	// UDP address creation for listening
	listenAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", listenIP, listenPort))
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// UDP connection for listening
	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		fmt.Println("Error creating UDP listener:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server listening on", listenAddr)

	// Buffer for incoming messages
	buffer := make([]byte, 1024)

	for {
		// Read from the connection
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			return
		}

		receivedMessage := string(buffer[:n])
		fmt.Printf("Received message '%s' from %s\n", receivedMessage, addr)
	}
}
