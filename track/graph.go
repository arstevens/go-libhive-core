package track

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	HistoryFile = "history.json"
)

// exchanges graphs should be created on the fly and deleted after use due to the
// amount of memory they take up in the processes memory
type ExchangeGraph struct {
	history map[string]float64
	parties []Party
}

func NewExchangeGraph(root string) (*ExchangeGraph, error) {
	var graph ExchangeGraph
	partyDirs := readDirectory(root)
	graph.parties = make([]Party, len(partyDirs))
	for i, partyDir := range partyDirs {
		graph.parties[i] = Party{id: filepath.Base(partyDir), fsLocation: partyDir}
	}

	hFile, err := os.Open(HistoryFile)
	if err != nil {
		return nil, err
	}
	defer hFile.Close()

	rawHistory, err := ioutil.ReadAll(hFile)
	if err != nil {
		return nil, err
	}

	graph.history = make(map[string]float64)
	err = json.Unmarshal(rawHistory, &graph.history)
	if err != nil {
		return nil, err
	}

	return &graph, nil
}

func (e ExchangeGraph) GetParty(id string) *Party {
	for _, party := range e.parties {
		if party.id == id {
			return &party
		}
	}

	return nil
}
