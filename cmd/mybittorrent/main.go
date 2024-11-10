package main

import (
	"bytes"
	"encoding/json"

	//"crypto/sha1"
	"fmt"
	"log"
	"os"
	"strings"

	bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

var _ = json.Marshal

func decodeBencode(bencodedString string) (interface{}, error) {
	reader := strings.NewReader(bencodedString)

	result, err := bencode.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("Couldn't decode the bencode string")
	}

	return result, nil
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
	var annouceUrl string
	var length int

	result, err := decodeBencode(bencodedString)

	jsonOutput, _ := json.Marshal(result) // no result here i.e null
	var data map[string]interface{}

	err = json.Unmarshal([]byte(jsonOutput), &data)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't unmarshal the JSON data")
	}

	annouceUrl = data["announce"].(string)

	if meta, ok := data["info"].(map[string]interface{}); ok {
		for key, value := range meta {
			if key == "length" {
				length = int(value.(float64))
			}
		}
	}

	if info, ok := data["info"]; ok {
		_, err := calculateInfoHash(info)
		if err != nil {
			return nil, nil, err
		}
	}

	return annouceUrl, length, nil
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

		if err != nil {
			log.Print("Couldn't read from stdin")
			os.Exit(1)
		}

		annonceUrl, length, err := extractTrackerURL(string(metaInfo))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Tracker URL:", annonceUrl)
		fmt.Println("Length:", length)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
