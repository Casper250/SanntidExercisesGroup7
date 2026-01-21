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
	defer close(out)

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

		case 0: //SLAVE

			fmt.Println("Backups start")

			addrRcv := net.UDPAddr{
				IP:   nil,
				Port: 20008,
			}

			connRcv, err := net.ListenUDP("udp", &addrRcv)
			if err != nil {
				log.Println("NEIIII", err)
			}

			buf := make([]byte, 1024)
			fmt.Println("Backups starts listening")
			for {

				connRcv.SetReadDeadline(time.Now().Add(3 * time.Second))
				n, _, err := connRcv.ReadFromUDP(buf)
				if err != nil {
					log.Println("Eg tek over skuta!")
					// Trigger a new master and backup
					masterInt = 1
					break

				} else {
					message := string(buf[:n])
					messageInt, err := strconv.Atoi(message)
					if err != nil {
						log.Println("Cannot convert to int", err)
					} else {
						counter = messageInt
						//fmt.Println("Received: ", counter)
					}
				}
			}
			connRcv.Close()

		case 1: // MASTER

			connSend, err := net.Dial("udp", "localhost:20008")
			if err != nil {
				log.Println("Nei nei nei nei")
			}

			out := make(chan int)
			go counting(counter, out)

			// Open a backup because the backup became master

			exec.Command("gnome-terminal", "--", "go", "run", "ProcessPairs.go", "0").Run()

			for c := range out {
				fmt.Println("Counter = ", c)

				// UDP Send a message on the PORT
				cString := strconv.Itoa(c)
				data := []byte(cString)

				_, err = connSend.Write(data)
				if err != nil {
					log.Println("Error sending message", err)
				}
			}
			connSend.Close()
		}

	}
}
