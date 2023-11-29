package main

import (
	"fmt"
	"net"
)

func calculateBroadcastAddress(ip net.IP, subnetMask net.IPMask) net.IP {
	// Invert the subnet mask
	invertedSubnetMask := net.IP(subnetMask)
	for i := range invertedSubnetMask {
		invertedSubnetMask[i] = ^invertedSubnetMask[i]
	}

	// Calculate the broadcast address
	broadcastAddress := make(net.IP, len(ip))
	for i := range ip {
		broadcastAddress[i] = ip[i] | invertedSubnetMask[i]
	}

	return broadcastAddress
}

func main() {
	// Get the local network interface and address
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var localIP net.IP
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
				if ip.To4() != nil {
					localIP = ip
					subnetMask = v.Mask
					break
				}
			}
		}
		if localIP != nil {
			break
		}
	}

	if localIP == nil {
		fmt.Println("Could not find a suitable network interface.")
		return
	}

	// Calculate the broadcast address
	broadcastIP := calculateBroadcastAddress(localIP, subnetMask)
	targetPort := 12345 // Arbitrary port

	// UDP address creation
	targetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcastIP.String(), targetPort))
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// UDP connection
	conn, err := net.DialUDP("udp", nil, targetAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		return
	}
	defer conn.Close()

	message := "Testing"
	messageBytes := []byte(message)

	// Send message
	_, err = conn.Write(messageBytes)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Printf("Sent message: %s to broadcast address %s\n", message, broadcastIP)
}
