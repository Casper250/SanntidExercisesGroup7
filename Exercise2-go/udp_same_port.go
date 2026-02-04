package udp_same_port

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

const (
	port = "20008"
	// Use "127.255.255.255" to test locally on one machine
	// Use "255.255.255.255" or your specific subnet (e.g. 192.168.1.255) for LAN
	broadcastAddr = "127.255.255.255"
)

func main() {
	senderID := fmt.Sprintf("node-%d", os.Getpid())
	counter := 0
	counter_str := "0"

	// UDP config stuff
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// 1. Allow multiple processes to use the port
				syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				// 2. Explicitly permit broadcasting from this socket
				syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
			})
		},
	}

	// Bind to all interfaces
	lp, err := lc.ListenPacket(context.Background(), "udp4", ":"+port)
	if err != nil {
		fmt.Printf("Listen error: %v\n", err)
		return
	}
	defer lp.Close()

	fmt.Printf("Instance %s active on port %s\n", senderID, port)

	// Receiver logic
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, addr, err := lp.ReadFrom(buffer)
			if err != nil {
				continue
			}
			msg := string(buffer[:n])
			fmt.Printf(" -> [RECV] from %s: %s\n", addr, msg)
		}
	}()

	// dst is a *UDPAddr 
	dst, _ := net.ResolveUDPAddr("udp4", broadcastAddr+":"+port)

	// Write a message to udp 
	ticker := time.NewTicker(10 * time.Millisecond)
	for range ticker.C {
		fmt.Printf("[SEND] broadcasting from %s...\n", senderID)

		counter_str = strconv.Itoa(counter)

		_, err := lp.WriteTo([]byte("Hello from "+senderID+" id = "+counter_str), dst)
		if err != nil {
			fmt.Printf("Send error: %v\n", err)
		}
		counter += 1
	}
}
