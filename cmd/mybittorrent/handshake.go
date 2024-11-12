package main

import (
	"fmt"
	"io"
)

type HandShakeStruct struct {
	ProtocolString string
	InfoHash       [20]byte
	PeerID         [20]byte
}

func (h *HandShakeStruct) NewHandshake(infoHash [20]byte, peerId [20]byte) *HandShakeStruct {
	return &HandShakeStruct{
		ProtocolString: "BitTorrent protocol",
		InfoHash:       infoHash,
		PeerID:         peerId,
	}
}

func (h *HandShakeStruct) Serialize() []byte {
	buf := make([]byte, len(h.ProtocolString)+49)
	buf[0] = byte(len(h.ProtocolString))
	curr := 1
	curr = copy(buf[curr:], h.ProtocolString)
	curr = copy(buf[curr:], make([]byte, 8))
	curr = copy(buf[curr:], h.InfoHash[:])
	curr = copy(buf[curr:], h.PeerID[:])

	return buf
}

func ReadHandshakeFromPeer(r io.Reader) (*HandShakeStruct, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	pstrlen := int(lengthBuf[0])

	if pstrlen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])

	h := HandShakeStruct{
		ProtocolString: string(handshakeBuf[0:pstrlen]),
		InfoHash:       infoHash,
		PeerID:         peerID,
	}

	return &h, nil
}
