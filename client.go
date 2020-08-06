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
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Starting up chat with remote computer...")
	inputChann := make(chan []byte)
	// recvChan := make(chan []byte)
	go startUDPPunching(inputChann)
	for {
		fmt.Print("--->")
		message, _ := reader.ReadString('\n')
		inputChann <- []byte(message)
	}
}

func startUDPPunching(inputChann chan []byte) {
	UDPConn, err := net.ListenUDP("udp", &net.UDPAddr{})

	if err != nil { // shouldn't ever happen
		fmt.Println("Issue opening udp socket.")
		panic(err)
	}
	clientAddress := resovleRemoteClientAddress(UDPConn)
	go sendToRemoteClient(UDPConn, clientAddress, inputChann) // punch hole in nat table
	go recieve(UDPConn)                                       // listen for incoming from opened port

}

func resovleRemoteClientAddress(UDPConn *net.UDPConn) *net.UDPAddr {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}

	for {
		data := make([]byte, 1024)
		fmt.Println(UDPConn.LocalAddr().String())
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

func sendToRemoteClient(UDPConn *net.UDPConn, remote *net.UDPAddr, inputChann chan []byte) {
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
