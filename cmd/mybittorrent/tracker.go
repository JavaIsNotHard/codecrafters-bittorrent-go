package main

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (torrent *TorrentFile) createURL(peerId [20]byte, port uint16) (string, error) {
	annouceUrl, err := url.Parse(torrent.Announce)

	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(torrent.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(torrent.Length)},
	}

	annouceUrl.RawQuery = params.Encode()
	return annouceUrl.String(), nil
}

func (torrent *TorrentFile) requestPeer(peerId [20]byte, port uint16) ([]Peer, error) {
	url, err := torrent.createURL(peerId, port)

	if err != nil {
		return []Peer{}, nil
	}

	resp, err := http.Get(url)

	if err != nil {
		return []Peer{}, nil
	}
	defer resp.Body.Close()

	trackerResponse := TrackerResponse{}

	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil {
		return []Peer{}, err
	}

	peers, err := Unmarshal([]byte(trackerResponse.Peers))
	if err != nil {
		return []Peer{}, nil
	}

	return peers, nil
}
