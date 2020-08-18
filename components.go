package main

import (
	"net"

	"github.com/google/uuid"
)

type peer struct {
	id   uuid.UUID
	addr net.UDPAddr
}

type peerPair struct {
	id       uuid.UUID
	internal net.UDPAddr
	external net.UDPAddr
}
