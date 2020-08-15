package main

import (
	"fmt"
	"net"
)

func startServer() {
	UDPListenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5000")
	conn, err := net.ListenUDP("udp", UDPListenAddr)
	if err != nil {
		fmt.Println("Couldn't open server at port 5000")
		panic(err)
	}

	clients := &clients{registered: make(map[string]peer)}
	for {
		buffer := make([]byte, 1024)
		_, addr, err := conn.ReadFromUDP(buffer)
		if err == nil { //message received succcessfully
			connectingPeer := peer{}
			if deserialize(buffer, &connectingPeer) == nil { //message format is correct
				connectingPeer.ExternalAddr = *addr
				go handleConnection(conn, clients, connectingPeer)
			}
		}
		fmt.Println("Recieved message from: ", addr.String())
	}
}

func handleConnection(serverSock *net.UDPConn, clients *clients, connectingPeer peer) {
	connectingPeerAddr := connectingPeer.ExternalAddr.String()
	if _, ok := clients.registered[connectingPeerAddr]; !ok {
		clients.registered[connectingPeerAddr] = connectingPeer // if is useless (might save lookup, or perhaps adds readability)
	}

	for address := range clients.registered {
		if address != connectingPeerAddr {
			serverSock.WriteToUDP(serialize(clients.registered[address]), &connectingPeer.ExternalAddr) // send any client which is not the connecting one
		}
	}
}
