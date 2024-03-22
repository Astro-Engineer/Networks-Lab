package main

import (
    "fmt"
    "net"
)

func main() {
    addr, err := net.ResolveUDPAddr("udp", ":12345")
    if err != nil {
        fmt.Println("Error resolving address:", err)
        return
    }
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("Error listening:", err)
        return
    }
    defer conn.Close()

    fmt.Println("UDP server started on port 12345")

    buffer := make([]byte, 1024)
    var total_msg = 0;
    for {
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("Error reading:", err)
            continue
        }
        msg_count ++;

        // Print the received message
        fmt.Printf("Received %d bytes from %s: %s\n", n, clientAddr, string(buffer[:n]))
        fmt.Printf("Total msgs %d\n", msg_count)
    }
}

