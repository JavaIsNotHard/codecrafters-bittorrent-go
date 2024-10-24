package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	reader := strings.NewReader(bencodedString)

	result, err := bencode.Decode(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't decode the bencode string")
	}

	jsonOutput, _ := json.Marshal(result)
	var data map[string]interface{}

	err = json.Unmarshal([]byte(jsonOutput), &data)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't decode the bencode string")
	}

	if meta, ok := data["info"].(map[string]interface{}); ok {
		for key, value := range meta {
			if key == "length" {
				length = value.(float64)
			}
		}
	}

	annouceUrl = data["announce"].(string)

	return annouceUrl, length, nil
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Errorf("Couldn't read from stdin")
			os.Exit(1)
		}

		// bencodedValue := os.Args[2]

		annonceUrl, length, err := extractTrackerURL(string(input))
		if err != nil {
			fmt.Println(err)
			return
		}

		// jsonOutput, _ := json.Marshal(annonceUrl)
		fmt.Println("Tracker URL:", annonceUrl)
		fmt.Println("Length:", length)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
