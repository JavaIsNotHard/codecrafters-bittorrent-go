package main

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	chokemsg        messageID = 0
	unchokedmsg     messageID = 1
	interestedmsg   messageID = 2
	uninterestedmsg messageID = 3
	havemsg         messageID = 4
	bitfieldmsg     messageID = 5
	requestmsg      messageID = 6
	piecemsg        messageID = 7
	cancelmsg       messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte // variable size bytes
}

func ReadMessageFromConn(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}

func (m *Message) SerializeMessage() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}
