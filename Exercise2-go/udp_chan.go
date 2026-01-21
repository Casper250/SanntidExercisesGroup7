package main

import (
	"fmt"
	"log"
	"net"
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

	chan_msg := make(chan string)
	chan_done := make(chan bool)

	go udp_rcvChan(config, chan_msg, chan_done)
	go udp_sendChan("Hei fra func 1", config, chan_msg, chan_done)
	go udp_sendChan("Hei fra func 2", config, chan_msg, chan_done)

	// Function for recieving calls to chan_done
	// Prevents deadlocking
	go func() {
		for {
			select {
			case <-chan_done:
				fmt.Println("Done signal detected")
				return
			default:
				// Some other work
				// Update stuff maybe
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	for {
		select {

		case msg := <-chan_msg:
			fmt.Println("Message received from channel: ", msg)

		case <-time.After(3 * time.Second):
			fmt.Println("No activity for 3 seconds ...")
			// Signal to recievers that no more data is coming
			close(chan_done)
			return
		}
	}

}

func udp_rcvChan(config ServerConfig, out chan<- string, done <-chan bool) {
	// address := config.ip + ":" + config.port_delimMsg

	addr := net.UDPAddr{
		//IP:   net.ParseIP(config.ip),
		IP:   nil,
		Port: 20001,
	}

	// Establish udp "connection"
	//conn, err := net.ListenUDP("udp", &addr)
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Get response from server
	buf := make([]byte, 1024)
	for {
		select {
		case <-done:
			fmt.Println("Receiver: Stopping early")
			return
		default:
			conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			out <- string(buf[:n])
		}
	}

}

func udp_sendChan(message string, config ServerConfig, out chan<- string, done <-chan bool) {

	address := config.ip + ":" + config.port_serverReceive

	// Establish udp "connection"
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Send 5 messages with a 1 second delay
	for i := 0; i < 5; i++ {
		data := []byte(message)
		_, err = conn.Write(data)
		if err != nil {
			log.Fatalln("Error sending message: ", err)
			continue
		}
		select {

		case <-done:
			fmt.Println("Sender: Stopping early ... ")
			return

		case <-time.After(1 * time.Second):
			// Next iteration

		}
	}
}
