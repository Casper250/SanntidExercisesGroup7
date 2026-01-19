package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

/*
Om udp
- oppnår ikke en connection slik som  med TCP
- sende og motta på forskjellige porter
- mest sannsynlig bruker vi UDP og lager ack funksjonene selv
-
*/

type SimulatorPorts struct {
	port_fixedSizeMsg string
	port_delimMsg     string
}

type ServerConfig struct {
	ip                 string
	port_serverReceive string
	port_serverReply   string
	port_fixedSizeMsg  string
	port_delimMsg      string
}

func main() {

	config := ServerConfig{
		ip:                 "172.26.161.38",
		port_serverReceive: "20000",
		port_serverReply:   "20001",
		port_fixedSizeMsg:  "34933",
		port_delimMsg:      "33546",
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go udp_send("Hei på deg", config, &wg)
	go udp_rcv(config, &wg)
	wg.Wait()

}

func udp_rcv(config ServerConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	// Establish udp "connection"
	//address := config.ip + ":" + config.port_delimMsg

	addr := net.UDPAddr{
		IP:   net.ParseIP(config.ip),
		Port: 20001,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Get response from server
	for {

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error reading from server: ", err)
			return
		}

		message := string(buf[:n])
		fmt.Println("Received: ", remoteAddr.IP, message)

	}

}

func udp_send(message string, config ServerConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	address := config.ip + ":" + config.port_serverReceive

	// Establish udp "connection"
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Send 10 messages
	i := 0
	for i < 10 {
		data := []byte(message)
		_, err = conn.Write(data)
		if err != nil {
			log.Fatalln("Error sending message", err)
		}
		fmt.Println("Sent: ", message)

		time.Sleep(100 * time.Millisecond)
		i++
	}

}
