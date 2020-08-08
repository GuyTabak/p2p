package main

import (
	"fmt"
	"net"
	"strconv"
)

// Clients ...
type Clients struct {
	registered map[string]bool // Need to change to IP:PORT
}

func startServer() {
	UDPListenAddr := &net.UDPAddr{Port: 5000}
	sock, err := net.ListenUDP("udp", UDPListenAddr)
	if err != nil {
		fmt.Println("Couldn't open server at port 5000")
		panic(err)
	}

	clients := &Clients{registered: make(map[string]bool)}
	for {
		buffer := make([]byte, 1024)
		_, addr, err := sock.ReadFromUDP(buffer)
		if err == nil {
			go handleConnection(addr, sock, clients)
		}
		fmt.Println("Recieved :", string(buffer))
	}

}

func handleConnection(connectingClient *net.UDPAddr, serverSock *net.UDPConn, clients *Clients) {
	stringAddress := connectingClient.String() + ":" + strconv.Itoa(connectingClient.Port) // ip:port
	if _, ok := clients.registered[stringAddress]; !ok {
		clients.registered[stringAddress] = true // if is useless (might save lookup, or perhaps adds readability)
	}

	for address := range clients.registered {
		if address != stringAddress {
			serverSock.WriteToUDP([]byte(address), connectingClient) // Send any client which is not the connecting one
			fmt.Printf("Debug:\nSent to client %v remote client %v", stringAddress, address)
		}
	}
}
