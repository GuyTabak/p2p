package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// StartClient is the entry point for any client
func StartClient() {
	inputChann := make(chan []byte)
	defer close(inputChann)

	fmt.Println("Starting up chat with remote computer...")
	go startUDPPunching(inputChann)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("--->")
		message, _ := reader.ReadString('\n')
		inputChann <- []byte(message)
	}
}

func startUDPPunching(inputChann chan []byte) {
	UDPConn, err := net.ListenUDP("udp", &net.UDPAddr{})

	if err != nil { //should never happen, sanity check
		fmt.Println("Issue opening udp socket.")
		panic(err)
	}
	clientAddress := resovleRemoteClientAddress(UDPConn)
	go send(UDPConn, clientAddress, inputChann)
	go recieve(UDPConn)

}

func resovleRemoteClientAddress(UDPConn *net.UDPConn) *net.UDPAddr {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}

	for {
		data := make([]byte, 1024)
		fmt.Println(UDPConn.LocalAddr().String()) // debug
		_, address, _ := UDPConn.ReadFromUDP(data)
		data = bytes.Trim(data, "\x00")
		if address.IP.String() == serverAddr.IP.String() {
			fmt.Println(string(data))
			UDPAddr, err := net.ResolveUDPAddr("udp", strings.TrimRight(string(data), "\n"))
			if err != nil {
				continue
			}
			return UDPAddr
		}
	}
}

func send(UDPConn *net.UDPConn, remote *net.UDPAddr, inputChann chan []byte) {
	go func() { //keep alive the udp punching
		empty := []byte{}
		for {
			UDPConn.WriteTo(empty, remote)
			time.Sleep(5)
		}
	}()

	for {
		message := <-inputChann
		fmt.Println("Sent following message to remote client: ", string(message))
		UDPConn.WriteTo(message, remote)
	}
}

func recieve(UDPConn *net.UDPConn) {
	for { // Receive loop
		buffer := make([]byte, 1024)
		_, addr, _ := UDPConn.ReadFromUDP(buffer)
		fmt.Printf("Recieved message: %v from remote %v.", string(buffer), addr.IP.String())
	}
}
