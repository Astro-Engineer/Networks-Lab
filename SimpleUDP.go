package main

import (
    "fmt"
    "net"
)

func main() {
    targetIP := "127.0.0.1" //local machine
    targetPort := 12345    //arbitrary port
    //UDP address creation
    targetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", targetIP, targetPort))
    if err != nil {
        fmt.Println("Error resolving address:", err)
        return
    }
    //UDP connection
    conn, err := net.DialUDP("udp", nil, targetAddr)
    if err != nil {
        fmt.Println("Error creating UDP connection:", err)
        return
    }
    //close connection (scheduled)
    defer conn.Close()
    message := "Testing"

    messageBytes := []byte(message)

    //send message
    _, err = conn.Write(messageBytes)
    if err != nil {
        fmt.Println("Error sending message:", err)
        return
    }
    fmt.Printf("Sent message: %s\n", message)
}
