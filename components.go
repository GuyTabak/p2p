package main

import "net"

type peer struct {
	InternalAddr net.UDPAddr
	ExternalAddr net.UDPAddr
}

type clients struct {
	registered map[string]peer
}
