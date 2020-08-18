package main

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

func startServer() {
	// Listen for clients registration
	UDPListenAddr, err := net.ResolveUDPAddr("udp", ":5000")
	conn, err := net.ListenUDP("udp", UDPListenAddr)
	if err != nil {
		fmt.Println("Couldn't open server at port 5000")
		panic(err)
	}

	peers := make(map[string]peerPair)
	for {
		buffer := make([]byte, 1024)
		_, addr, err := conn.ReadFromUDP(buffer)

		if err == nil { //message received succcessfully
			internalPeer := peer{}
			if deserialize(buffer, &internalPeer) == nil { //message format is correct
				go handleConnection(conn, peers, internalPeer, addr)
			}
		}
		fmt.Println("Recieved message from: ", addr.String())
	}
}

func handleConnection(serverConn *net.UDPConn, peers map[string]peerPair, internal peer, connectingAddr *net.UDPAddr) {
	connectingPeerAddr := connectingAddr.String()

	if _, ok := peers[connectingPeerAddr]; !ok {
		id := uuid.New()
		internal.id = id
		peers[connectingPeerAddr] = peerPair{id: id, internal: internal.addr, external: *connectingAddr}
	}

	for address, pair := range peers {
		if address != connectingPeerAddr {
			externalPeerInfo := NetworkMessage{messageType: update, payload: serialize(peer{id: pair.id, addr: pair.external})}
			internalPeerInfo := NetworkMessage{messageType: update, payload: serialize(peer{id: pair.id, addr: pair.internal})}

			serverConn.WriteToUDP(serialize(externalPeerInfo), connectingAddr)
			serverConn.WriteToUDP(serialize(internalPeerInfo), connectingAddr) // send any client which is not the connecting one
		}
	}
}
