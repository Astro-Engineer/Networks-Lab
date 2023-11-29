package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"strconv"
)

func calculateBroadcastAddress(ip net.IP, subnetMask net.IPMask) net.IP {
	// Ensure the IP and subnetMask are IPv4
	ip = ip.To4()
	if ip == nil {
		fmt.Println("Invalid IPv4 address.")
		return nil
	}

	// Invert the subnet mask
	invertedSubnetMask := net.IP(subnetMask)
	for i := range invertedSubnetMask {
		invertedSubnetMask[i] = ^invertedSubnetMask[i]
	}

	// Calculate the broadcast address
	if len(ip) == len(invertedSubnetMask) {
		broadcastAddress := make(net.IP, len(ip))
		for i := range ip {
			broadcastAddress[i] = ip[i] | invertedSubnetMask[i]
		}
		return broadcastAddress
	}

	fmt.Println("Invalid IPv4 address or subnet mask lengths.")
	return nil
}

func main() {
	// Set Priority
	//var priority int
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./program <priority>")
		return
	}

	// Parse the prio from command-line arguments
	priority, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error parsing target port:", err)
		return
	}
	
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

	// Channel to handle termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to handle incoming messages
	go func() {
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
	}()

	// Goroutine to handle user input and send messages
	go func() {
		// Get the local network interface and IPv4 address
		interfaces, err := net.Interfaces()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var localIPv4 net.IP
		var subnetMask net.IPMask

		for _, iface := range interfaces {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					ip := v.IP
					if ip.To4() != nil && !ip.IsLoopback() {
						localIPv4 = ip
						subnetMask = v.Mask
						break
					}
				}
			}
			if localIPv4 != nil {
				break
			}
		}

		if localIPv4 == nil {
			fmt.Println("Could not find a suitable non-loopback IPv4 network interface.")
			return
		}

		// Print the local IPv4 address
		fmt.Printf("IPv4 broadcast address: %s\n", localIPv4)
		// Print the subnet mask
		fmt.Printf("Subnet Mask: %s\n", subnetMask)

		// Calculate the broadcast address
		broadcastIP := calculateBroadcastAddress(localIPv4, subnetMask)
		targetPort := 12345 // Arbitrary port

		// Print the target IPv4 address
		fmt.Printf("Sending to IPv4 broadcast address: %s\n", broadcastIP)

		// UDP address creation
		targetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcastIP.String(), targetPort))
		if err != nil {
			fmt.Println("Error resolving address:", err)
			return
		}

		// UDP connection
		sendConn, err := net.DialUDP("udp", nil, targetAddr)
		if err != nil {
			fmt.Println("Error creating UDP connection:", err)
			return
		}
		defer sendConn.Close()

		// Start a goroutine to listen for user input
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			//change this behavior to periodically (we are using the scanner so we can control when the node sends msgs)
			message := scanner.Text() + strconv.Itoa(priority)
			messageBytes := []byte(message)

			// Check if the user input is not empty
			if strings.TrimSpace(message) != "" {
				// Send message
				_, err := sendConn.Write(messageBytes)
				if err != nil {
					fmt.Println("Error sending message:", err)
					return
				}
			}
		}
	}()

	// Wait for termination signals
	<-signalCh
	fmt.Println("\nServer shutting down.")
}
