package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded,  err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
