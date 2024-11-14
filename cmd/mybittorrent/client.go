package main

import (
	"net"
)

type Client struct {
	Conn     net.Conn
	IsChoked bool
	Bitfield Bitfield
	peer     Peer
	infoHash [20]byte
	peerID   [20]byte
}
