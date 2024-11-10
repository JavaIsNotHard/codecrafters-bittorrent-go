package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"

	//"crypto/sha1"
	"fmt"
	"log"
	"os"
	"strings"

	bencode "github.com/jackpal/bencode-go"
)

var _ = json.Marshal

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

func generateInfoHash(info *bencodeInfo) (interface{}, error) {
	jsonBencodedData, err := json.Marshal(info)

	if err != nil {
		return nil, err
	}

	hash := sha1.New()
	hash.Write([]byte(jsonBencodedData))
	return hex.EncodeToString(hash.Sum(nil)), nil
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

func extractPieces(info bencodeInfo) ([][20]byte, error) {
	hashlen := 20
	buf := []byte(info.Pieces)
	if len(buf)%hashlen != 0 {
		err := fmt.Errorf("received malformed pieces of length: %d", len(buf))
		return nil, err
	}

	numhashes := len(buf)
	hashes := make([][20]byte, numhashes)

	for i := 0; i < numhashes; i++ {
		copy(hashes[i][:], buf[i*hashlen:(i+1)*hashlen])
	}

	return hashes, nil
}

func extractTrackerURL(bencodedString string) (interface{}, *bencodeInfo, error) {
	result, err := decodeBencode(bencodedString)

	if err != nil {
		return nil, nil, err
	}

	return result.Announce, &result.Info, nil
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			log.Print("Couldn't decode the string")
			os.Exit(1)
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else if command == "info" {
		fileName := os.Args[2]

		metaInfo, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal("Couldn't open file", fileName)
			os.Exit(1)
		}

		annonceUrl, info, err := extractTrackerURL(string(metaInfo))
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		// infoHash, err := generateInfoHash(info)
		// if err != nil {
		// 	log.Print(err)
		// 	os.Exit(1)
		// }

		// pieces, err := extractPieces(*info)
		// if err != nil {
		// 	log.Print(err)
		// 	os.Exit(1)
		// }

		// fmt.Println(pieces)

		fmt.Println("Tracker URL:", annonceUrl)
		fmt.Println("Length:", info.Length)
		//fmt.Println("Info Hash:", infoHash)
		//fmt.Println("Pieces Length:", info.PieceLength)
		//fmt.Println()
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
