package main

import (
	"net"

	"github.com/google/uuid"
)

type peer struct {
	id   uuid.UUID
	addr net.UDPAddr
}
