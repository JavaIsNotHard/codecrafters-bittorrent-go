package main

import (
	"bytes"
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

func decodeBencode(bencodedString string) (*bencodeTorrent, error) {
	reader := strings.NewReader(bencodedString)

	bto := bencodeTorrent{}
	err := bencode.Unmarshal(reader, &bto)
	if err != nil {
		return nil, err
	}

	return &bto, nil
}

func calculateInfoHash(infoMap interface{}) (string, error) {
	newInfoMap := infoMap.(map[string]interface{})
	jsonOutput, _ := json.Marshal(newInfoMap)

	fmt.Println(string(jsonOutput))

	var buffer bytes.Buffer

	_, err := buffer.Write(jsonOutput)
	if err != nil {
		return "", err
	}

	err = bencode.Marshal(&buffer, infoMap)
	if err != nil {
		return "", err
	}

	fmt.Println("The bencode value is: ", buffer.String())

	return "", nil
}

func extractTrackerURL(bencodedString string) (interface{}, interface{}, error) {
	result, err := decodeBencode(bencodedString)

	if err != nil {
		return nil, nil, err
	}

	return result.Announce, result.Info.Length, nil
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

		annonceUrl, length, err := extractTrackerURL(string(metaInfo))
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		fmt.Println("Tracker URL:", annonceUrl)
		fmt.Println("Length:", length)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
