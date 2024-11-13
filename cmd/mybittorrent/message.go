package main

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
