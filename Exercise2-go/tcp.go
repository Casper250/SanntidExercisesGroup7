package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go tcp_send("Hello from 1", &wg)
	go tcp_send("Hello from 2", &wg)
	go tcp_send("Hello from 3", &wg)
	wg.Wait()
}

func tcp_send(message string, wg *sync.WaitGroup) {
	defer wg.Done()
	var terminatorChar = "\x00"
	var counter = 0
	data := []byte(message + terminatorChar)

	// TCP socket
	// Portnummeret kom på VsCode med den lokale severen
	// var port_delimMessage = "33546"
	// address format: a.b.c.d:port
	var port_fixedMessageSize = "34933"
	var ip_serverSide = "172.26.161.38"
	var address = ip_serverSide + ":" + port_fixedMessageSize

	// Make a tcp connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln(conn, "Connection failed")
		return
	}
	defer conn.Close()

	for counter < 5 {

		// Send a message
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Read a response
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
		}

		// Print response
		fmt.Println("DATA FROM SERVER: ", string(buf[:n]))

		time.Sleep(333 * time.Millisecond)
		counter++
	}

}
