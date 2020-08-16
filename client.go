package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func startClient() {
	inputChann := make(chan []byte)
	defer close(inputChann)

	fmt.Println("Starting up chat with remote computer...")
	go startUDPPunching(inputChann)

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		inputChann <- []byte(message)
	}
}

func startUDPPunching(inputChann chan []byte) {
	UDPConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0}) //bind socket to random port, on udp layer
	fmt.Println("My ip: ", UDPConn.LocalAddr())

	if err != nil { //should never happen, sanity check
		panic(err)
	}
	remotePeer := resovleRemoteClientAddress(UDPConn)
	go send(UDPConn, remotePeer, inputChann)
	go recieve(UDPConn)

}

func resovleRemoteClientAddress(UDPConn *net.UDPConn) peer {
	UDPIP := getOutboundIP().IP.String()
	UDPPort := strconv.Itoa(UDPConn.LocalAddr().(*net.UDPAddr).Port)

	UDPAddr, _ := net.ResolveUDPAddr("udp", UDPIP+":"+UDPPort)
	encodedAddr := serialize(peer{InternalAddr: *UDPAddr})
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000") // TODO: config
	if err != nil {
		panic(err)
	}

	defer UDPConn.SetReadDeadline(time.Time{})
	for {
		_, err := UDPConn.WriteTo(encodedAddr, serverAddr)

		UDPConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		data := make([]byte, 1024)
		len, address, _ := UDPConn.ReadFromUDP(data)

		if err != nil {
			continue // wait for server response
		}

		if address != nil && len > 0 && address.IP.String() == serverAddr.IP.String() {
			var remotePeer peer
			deserialize(data, &remotePeer)

			return remotePeer
		}
	}
}

func send(UDPConn *net.UDPConn, remote peer, inputChann chan []byte) {

	// go func() { //keep alive
	// 	for {
	// 		UDPConn.WriteTo([]byte{}, &remote.InternalAddr)
	// 		UDPConn.WriteTo([]byte{}, &remote.ExternalAddr)
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	for {
		message := <-inputChann
		UDPConn.WriteTo(message, &remote.InternalAddr)
		UDPConn.WriteTo(message, &remote.ExternalAddr)
	}
}

// TODO: Avoid print of managment messages from server (might occure upon slow response)
func recieve(UDPConn *net.UDPConn) {
	for { // Receive loop
		buffer := make([]byte, 1024)
		len, addr, err := UDPConn.ReadFromUDP(buffer)
		if err == nil {
			if len > 0 { //avoid keepalive
				fmt.Printf("From %v: %v", addr.String(), string(buffer))
			}
		} else {
			fmt.Println(err)
		}
	}
}
