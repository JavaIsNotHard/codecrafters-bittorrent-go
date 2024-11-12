package main

// perform handshake will all the peers using goroutines
// chocked -> no data will be sent until unchoking happens
//
// whenever one side is interested and the other side is not chocked
// downloaders should keep several pieces request queued up at once in order to get good performance (known as pipelining)
// request which cannot be written out to TCP buffer should immediately be queued in memory rather than in application-level network buffer
// the handshake starts with character nineteen (decimal) followed by the string 'BitTorrent Protocol'
// all later integers sent in the protocol are encoded as four bytes big-endian
// after the fixed header comes eight reserved bytes whose value are all zeros
// next is the sha1 hash of the info dict from the bencoded torrent file

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net"
	"time"

	//"crypto/sha1"
	"fmt"
	"log"
	"os"
	"strings"

	bencode "github.com/jackpal/bencode-go"
)

const Port uint16 = 6881

var _ = json.Marshal

type TorrentFile struct {
	Announce     string
	InfoHash     [20]byte
	PiecesHashes [][20]byte
	PieceLength  int
	Length       int
	Name         string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (i *bencodeInfo) generateInfoHash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)

	if err != nil {
		return [20]byte{}, nil
	}

	hash := sha1.Sum(buf.Bytes())

	return hash, nil
}

func decodeBencode(bencodedString string) (*bencodeTorrent, error) {
	reader := strings.NewReader(bencodedString)

	bto := bencodeTorrent{}
	err := bencode.Unmarshal(reader, &bto)
	if err != nil {
		return nil, err
	}

	return &bto, nil
}

func (i *bencodeInfo) extractPieces() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Pieces)
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (bt *bencodeTorrent) toTorrent() (TorrentFile, error) {
	infoHash, err := bt.Info.generateInfoHash()

	if err != nil {
		return TorrentFile{}, nil
	}

	piecesHashes, err := bt.Info.extractPieces()
	if err != nil {
		return TorrentFile{}, nil
	}

	return TorrentFile{
		Announce:     bt.Announce,
		InfoHash:     infoHash,
		PiecesHashes: piecesHashes,
		PieceLength:  bt.Info.PieceLength,
		Length:       bt.Info.Length,
		Name:         bt.Info.Name,
	}, nil
}

func (torrent *TorrentFile) printHashList() {
	for i := 0; i < len(torrent.PiecesHashes); i++ {
		fmt.Println(hex.EncodeToString(torrent.PiecesHashes[i][:]))
	}
}

func generatePeerID() [20]byte {
	peerID := [20]byte{}
	copy(peerID[:], "-GO0001-"+"123456789012")
	return peerID
}

const protocolName = "BitTorrent protocol"

func (torrentData *Torrent) createConnection(address string) error {
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	var buffer bytes.Buffer
	buffer.WriteByte(byte(len(protocolName)))
	buffer.WriteString(protocolName)
	buffer.Write(make([]byte, 8))
	buffer.Write(torrentData.InfoHash[:])
	fmt.Println(hex.EncodeToString(torrentData.InfoHash[:]))
	buffer.Write(torrentData.PeerID[:])

	count, err := conn.Write(buffer.Bytes())
	fmt.Println(count)

	resp, err := ReadHandshakeFromPeer(conn)
	fmt.Println(count)
	if err != nil {
		if err == os.ErrClosed {
			fmt.Println("Connection closed by server.")
		}
		if err.Error() == "EOF" {
			fmt.Println("Reached end of file.")
		}
		fmt.Println("Error reading from connection:", err)
		fmt.Println("failed:", err)
	}

	fmt.Println(hex.EncodeToString(resp.InfoHash[:]))
	return nil
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := bencode.Decode(strings.NewReader(bencodedValue))
		if err != nil {
			log.Print("Couldn't decode the string")
			os.Exit(1)
		}

		jsonOutput, err := json.Marshal(decoded)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(jsonOutput))
	} else if command == "info" {
		fileName := os.Args[2]

		torrentContent, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal("Couldn't open file", fileName)
			os.Exit(1)
		}

		decodedMetaInfo, err := decodeBencode(string(torrentContent))
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		torrent, err := decodedMetaInfo.toTorrent()

		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		fmt.Println("Tracker URL:", torrent.Announce)
		fmt.Println("Length:", torrent.Length)
		fmt.Println("Info Hash:", hex.EncodeToString(torrent.InfoHash[:]))
		fmt.Println("Piece Length:", torrent.PieceLength)
		fmt.Println("Pieces Hashes:")
		torrent.printHashList()

	} else if command == "peers" {
		fileName := os.Args[2]

		torrentContent, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal("Couldn't open file", fileName)
			os.Exit(1)
		}

		decodedMetaInfo, err := decodeBencode(string(torrentContent))
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		torrent, err := decodedMetaInfo.toTorrent()
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		var peerId [20]byte

		_, err = rand.Read(peerId[:])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		peers, err := torrent.requestPeer(peerId, Port)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		for _, value := range peers {
			fmt.Println(value)
		}

		// torrentdata := Torrent{
		// 	Peers:       peers,
		// 	PeerID:      peerId,
		// 	InfoHash:    torrent.InfoHash,
		// 	PiecesHash:  torrent.PiecesHashes,
		// 	PieceLength: torrent.PieceLength,
		// 	Length:      torrent.Length,
		// 	Name:        torrent.Name,
		// }

		// fmt.Println(peers[0].String())
		// fmt.Println(hex.EncodeToString(torrentdata.PeerID[:]))
		// torrentdata.createConnection(peers[0].String())

		// fmt.Println("All peers connected")

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
