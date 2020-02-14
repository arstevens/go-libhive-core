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
	parties []*Party
}

func NewExchangeGraph(root string) (*ExchangeGraph, error) {
	var graph ExchangeGraph
	partyDirs := readDirectory(root)
	graph.parties = make([]*Party, len(partyDirs))
	for i, partyDir := range partyDirs {
		graph.parties[i] = &Party{id: filepath.Base(partyDir), fsLocation: partyDir}
	}

	historyPath := root + "/" + HistoryFile
	var hFile *os.File
	var err error
	if !fileExists(HistoryFile) {
		hFile, err = os.Create(historyPath)
		if err != nil {
			return nil, err
		}
	} else {
		hFile, err = os.Open(root + "/" + HistoryFile)
		if err != nil {
			return nil, err
		}
	}
	defer hFile.Close()

	// Read starting values for each party
	rawHistory, err := ioutil.ReadAll(hFile)
	if err != nil {
		return nil, err
	}
	history := make(map[string]float64)
	err = json.Unmarshal(rawHistory, &history)
	if err != nil {
		return nil, err
	}

	// Store starting values in Party objects
	for _, party := range graph.parties {
		party.history = history[party.id]
	}

	return &graph, nil
}

func (e *ExchangeGraph) AddParty(p *Party) {
	e.parties = append(e.parties, p)
}

func (e *ExchangeGraph) GetParty(id string) *Party {
	for _, party := range e.parties {
		if party.id == id {
			return party
		}
	}

	return nil
}

func (e *ExchangeGraph) Compress(outpath string) error {
	// Compress transaction histories into single values
	newHistory := make(map[string]float64)
	for _, party := range e.parties {
		sum, err := party.SumTransactions()
		if err != nil {
			return err
		}
		newHistory[party.id] = sum
	}

	// Prepare compressed data and write to disk
	rawHistory, err := json.Marshal(newHistory)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outpath, rawHistory, 0466) // r--rw-rw-
}

// Clears all transactions in root directory. Ignores history file
func (e *ExchangeGraph) DeleteTransactions() error {
	for _, party := range e.parties {
		err := os.Remove(party.fsLocation)
		if err != nil {
			return err
		}
	}
	return nil
}
