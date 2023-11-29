package main

import (
	"fmt"
	"net"
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
				if ip.To4() != nil && ip.IsLoopback() == false {
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

	fmt.Printf("Sent message: %s\n", message)
}

