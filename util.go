package main

import (
	"net"
)

func getOutboundIP() net.UDPAddr {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	return *conn.LocalAddr().(*net.UDPAddr)
}
