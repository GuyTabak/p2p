package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

func startClient() {
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
	UDPConn, err := net.ListenUDP("udp", &net.UDPAddr{}) //bind socket to random port, on udp layer

	if err != nil { //should never happen, sanity check
		fmt.Println("Issue opening udp socket.")
		panic(err)
	}
	clientAddress := resovleRemoteClientAddress(UDPConn)
	go send(UDPConn, clientAddress, inputChann)
	go recieve(UDPConn)

}

func resovleRemoteClientAddress(UDPConn *net.UDPConn) *net.UDPAddr {
	serverAddr, err := net.ResolveUDPAddr("udp", "3.22.98.195:5000") // UPDATE SERVER IP HERE
	if err != nil {
		panic(err)
	}

	for { //send 'keep alive' wait loop
		data := make([]byte, 1024)
		UDPConn.WriteTo([]byte("Register request."), serverAddr)
		UDPConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, address, _ := UDPConn.ReadFromUDP(data)
		if err != nil {
			continue // wait for server response
		}

		if address != nil && address.IP.String() == serverAddr.IP.String() {
			UDPAddr, err := net.ResolveUDPAddr("udp", string(bytes.Trim(data, "\x00")))
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
