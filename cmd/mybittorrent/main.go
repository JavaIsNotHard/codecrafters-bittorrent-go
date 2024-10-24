package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func extractTrackerURL(bencodedString string) (interface{}, interface{}, error) {
	var annouceUrl string
	var length float64

    result, err := decodeBencode(bencodedString)

	jsonOutput, _ := json.Marshal(result) // no result here i.e null
    fmt.Println(string(jsonOutput))
	var data map[string]interface{}

	err = json.Unmarshal([]byte(jsonOutput), &data)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't unmarshal the JSON data")
	}

	annouceUrl = data["announce"].(string)

	if meta, ok := data["info"].(map[string]interface{}); ok {
		for key, value := range meta {
			if key == "length" {
				length = value.(float64)
			}
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
        
        fmt.Println(string(metaInfo))

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
