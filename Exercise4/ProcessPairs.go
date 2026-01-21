package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func counting(start int, out chan<- int) {

	counter := start
	fmt.Println("Now this process starts to count")

	for {
		time.Sleep(time.Second) //Later on should be more careful with the seconds
		counter += 1
		//fmt.Println("value: ", counter)
		out <- counter

		if counter == 100 {
			break
		}
	}

}

func main() {

	masterInt, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Println("Master ikke int")
	}

	fmt.Println("Master: ", masterInt)
	counter := 0

	for {
		switch masterInt {
		//SLAVE
		case 0:
			// UDP recieve/ reset timer

			//timer := 0
			fmt.Println("Backups start")

			addrRcv := net.UDPAddr{
				IP:   nil,
				Port: 20008,
			}

			connRcv, err := net.ListenUDP("udp", &addrRcv)
			if err != nil {
				log.Fatalln("NEIIII", err)
			}

			buf := make([]byte, 1024)
			fmt.Println("Backups starts listening")
			for {

				connRcv.SetReadDeadline(time.Now().Add(2 * time.Second))
				n, _, err := connRcv.ReadFromUDP(buf)
				if err != nil {
					log.Println("Eg tek over skuta!")
					masterInt = 1
					connRcv.Close()
					break

				} else {
					message := string(buf[:n])
					messageInt, err := strconv.Atoi(message)
					if err != nil {
						log.Println("Cannot convert to int 7", err)
					} else {
						counter = messageInt
						//fmt.Println("Received: ", counter)
					}
				}

			}
		// MASTER
		case 1:

			connSend, err := net.Dial("udp", "localhost:20008")
			if err != nil {
				log.Fatalln("Nei nei nei nei")
			}
			defer connSend.Close()

			out := make(chan int)

			go counting(counter, out)

			// Open a backup because the backup became master
			//fmt.Println("Master want to initialize backup")
			exec.Command("gnome-terminal", "--", "go", "run", "ProcessPairs.go", "0").Run()

			for {

				select {
				case c := <-out:
					fmt.Println("Counter = ", c)

					// UDP Send a message on the PORT
					cString := strconv.Itoa(c)
					data := []byte(cString)

					_, err = connSend.Write(data)
					if err != nil {
						log.Fatalln("Error sending message", err)
					}
				}
			}

		}
	}

}
