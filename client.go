package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type manager struct {
	myPeerID   uuid.UUID
	registered map[uuid.UUID]peer   // peers from which i received ack
	newPeers   map[uuid.UUID][]peer // peers yet to return ack

	broadcast      chan []byte
	shutdownSignal bool
	lock           *sync.Mutex
}

func startClient() {
	UDPConn, err := net.ListenUDP("udp", &net.UDPAddr{}) //bind socket to random port, on udp layer

	//should never happen, sanity check
	if err != nil {
		panic(err)
	}

	manager := &manager{}
	go listenLoop(UDPConn, manager)
	go manager.maintain(UDPConn)
	register(UDPConn)

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		manager.broadcast <- []byte(message)
	}

}

// TODO: Consider case where server didn't receive registration -> retry? tcp?
func register(UDPConn *net.UDPConn) {
	// send client's internal ip, for in-lan comuunication
	UDPIP := getOutboundIP().IP.String()
	UDPPort := strconv.Itoa(UDPConn.LocalAddr().(*net.UDPAddr).Port)
	UDPAddr, _ := net.ResolveUDPAddr("udp", UDPIP+":"+UDPPort)
	data := serialize(peer{addr: *UDPAddr})

	serverAddr, err := net.ResolveUDPAddr("udp", "3.22.98.195:5000") // config
	if err != nil {
		panic(err)
	}

	UDPConn.WriteToUDP(data, serverAddr)
}

// start of message RECIEVE handling
func listenLoop(UDPConn *net.UDPConn, manager *manager) {
	for {
		// read data from listening udp connection
		data := make([]byte, 4096)
		len, addr, err := UDPConn.ReadFromUDP(data)
		if len == 0 || err != nil {
			continue
		}

		// deserialize data into message struct
		msg := &NetworkMessage{}
		err = deserialize(data, msg)
		msg.addr = *addr // the reciever is responsible of popuplating the address field, as the sender isn't aware of it
		if err != nil {
			continue
		}

		// handle message
		switch msg.messageType {
		case message:
			fmt.Printf("From %v: %s", addr.String(), msg.payload) //  assure one can convert []byte to string with %s
		case handshake:
			handleHandshake(*msg, UDPConn, manager)
		case keepalive: // add log or use it for garbage collector
		case update:
			manager.update(*msg)
		case close:
			manager.close(*msg)
		default:
		}
	}
}

// TODO: Test
func handleHandshake(msg NetworkMessage, conn *net.UDPConn, m *manager) {
	switch msgType := binary.LittleEndian.Uint16(msg.payload); msgType {
	case synack: // peer sent ack
		//TODO: move peer from new to registered
	case syn: // peer requires ack
		var data []byte
		binary.LittleEndian.PutUint16(data, synack)

		msg = NetworkMessage{messageType: handshake, peerID: m.myPeerID, payload: data}
		conn.WriteToUDP(serialize(msg), &msg.addr)
	}
}

// TODO: Test
func (m *manager) markRegistered(p peer) {
	if _, ok := m.registered[p.id]; !ok {
		m.registered[p.id] = p
	}

	if peers, ok := m.newPeers[p.id]; ok {
		for index, possiblePeer := range peers {
			if possiblePeer.addr.String() == p.addr.String() {
				peers[index] = peers[len(peers)-1]
				peers = peers[:len(peers)-1] //TODO: Verify that peers is indeed pointer to the connecting one.
			}
		}
	}
}

// TODO: Test
func (m *manager) update(msg NetworkMessage) {
	newPeer := &peer{}
	err := deserialize(msg.payload, newPeer)
	if err != nil {
		return
	}

	if peers, ok := m.newPeers[msg.peerID]; ok {
		for _, peer := range peers {
			if peer.addr.String() == newPeer.addr.String() {
				return
			}
		}
	}

	m.newPeers[msg.peerID] = []peer{*newPeer}
}

func (m *manager) close(msg NetworkMessage) {
	// The connecting items will be cleaned manager's garbage collector
	delete(m.registered, msg.peerID)
}

// end of message RECIEVE handling

/*
	Responisble of the main logic:
	- Send keep alives to peers (established peers)
	- Attemps connecting the new peers
	- removes unreachable peers
*/
func (m *manager) maintain(UDPConn *net.UDPConn) {
	go sendLoop(m, UDPConn)
	for !m.shutdownSignal {
		go sendKeepAlive(m, UDPConn)
		go sendSYN(m, UDPConn)
		time.Sleep(10 * time.Second)
	}
}

func sendKeepAlive(m *manager, UDPConn *net.UDPConn) {
	// Check how to send multicast.
	for !m.shutdownSignal {
		message := serialize(NetworkMessage{messageType: keepalive})
		for _, peer := range m.registered {
			UDPConn.WriteToUDP(message, &peer.addr)
		}
	}
}

func sendSYN(m *manager, UDPConn *net.UDPConn) {
	var data []byte
	binary.LittleEndian.PutUint16(data, syn)

	SYNMessage := NetworkMessage{peerID: m.myPeerID, messageType: handshake, payload: data}

	for _, peers := range m.newPeers {
		for _, peer := range peers {
			UDPConn.WriteToUDP(serialize(SYNMessage), &peer.addr)
		}
	}
}

func sendLoop(m *manager, UDPConn *net.UDPConn) {
	// check broadcast/multiast
	for !m.shutdownSignal {
		select {
		case data := <-m.broadcast:
			msg := &NetworkMessage{messageType: message, payload: data, peerID: m.myPeerID}
			for _, peer := range m.registered {
				UDPConn.WriteToUDP(serialize(msg), &peer.addr)
			}
		}
	}
}
