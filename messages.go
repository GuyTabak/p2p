package main

import (
	"net"

	"github.com/google/uuid"
)

// NetworkMessage ...
type NetworkMessage struct {
	messageType uint16
	payload     []byte
	peerID      uuid.UUID
	addr        net.UDPAddr //sender's addr, optiona
}

const (
	message   uint16 = 0
	stream    uint16 = 1
	handshake uint16 = 2
	keepalive uint16 = 3
	update    uint16 = 4
	close     uint16 = 5
)

const (
	syn    uint16 = 0
	synack uint16 = 1
)

//TODO: Read regarding uint16 and why there is no binary.put for it
