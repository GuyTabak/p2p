package main

import (
	"fmt"
	"net"
)

// Clients ...
type Clients struct {
	registered map[string]*net.UDPAddr // Need to change to IP:PORT
}

func startServer() {
	UDPListenAddr := &net.UDPAddr{Port: 5000}
	sock, err := net.ListenUDP("udp", UDPListenAddr)
	if err != nil {
		fmt.Println("Couldn't open server at port 5000")
		panic(err)
	}

	clients := &Clients{registered: make(map[string]*net.UDPAddr)}
	for {
		buffer := make([]byte, 1024)
		_, addr, err := sock.ReadFromUDP(buffer)
		if err == nil {
			go handleConnection(addr, sock, clients)
		}
		fmt.Println("Recieved :", string(buffer))
	}

}

func handleConnection(addr *net.UDPAddr, serverSock *net.UDPConn, clients *Clients) {
	if _, ok := clients.registered[addr.IP.String()]; !ok {
		clients.registered[addr.IP.String()+":"+string(addr.Port)] = addr
	}

	for addrStr, addr := range clients.registered {
		if addrStr != addr.IP.String()+":"+string(addr.Port) {
			serverSock.WriteToUDP([]byte(addrStr), addr)
			fmt.Printf("Debug:\nSent to cleint %v remote client %v", addr, addr)
		} else {
			fmt.Printf("Debug.")
		}
	}
}
