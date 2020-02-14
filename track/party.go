package track

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Is there a way to have nodes only keep info on selective parties &
// still enable synchronous consensus?
type Party struct {
	id         string
	fsLocation string
	history    float64
}

func NewParty(id string, loc string, hist float64) *Party {
	fsLoc := loc + "/" + id
	if !fileExists(fsLoc) {
		os.Mkdir(fsLoc, 0777)
	}
	return &Party{id: id, fsLocation: fsLoc, history: hist}
}

func (p *Party) Id() string {
	return p.id
}

func (p *Party) AddTransaction(t Transaction) error {
	// Add file to the folder for this parties transactions
	tFile := p.fsLocation + "/" + t.Id()
	if fileExists(tFile) {
		return errors.New("Transaction already exists. Transactions are immutable")
	}

	transactionFile, err := os.Create(tFile)
	if err != nil {
		return err
	}
	defer transactionFile.Close()

	serial := t.Marshal()
	_, err = transactionFile.Write(serial)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Link this transaction to the other members of this transaction
	parties := t.Parties()
	parentDir := filepath.Dir(p.fsLocation)

	for _, party := range parties {
		if party != p.Id() {
			partyFound := fileExists(parentDir + "/" + party)
			if partyFound {
				symPath := parentDir + "/" + party + "/" + t.Id()
				fmt.Println(transactionFile.Name() + " : " + symPath)
				err = os.Symlink(transactionFile.Name(), symPath)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}

	return err
}

func (p *Party) SumTransactions() (float64, error) {
	transactionPaths := readDirectory(p.fsLocation)
	transcations, err := parseTransactions(transactionPaths)
	if err != nil {
		return -1.0, err
	}

	sum := p.history
	for _, curTransaction := range transcations {
		sum += (*curTransaction).GetAmountExchanged(p.id)
	}

	return sum, nil
}

func readDirectory(root string) []string {
	dirPaths := make([]string, 0)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		dirPaths = append(dirPaths, path)
		return nil
	})

	return dirPaths
}

func parseTransactions(transactionPaths []string) ([]*Transaction, error) {
	transactions := make([]*Transaction, len(transactionPaths))
	for _, fpath := range transactionPaths {
		file, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}

		curTransaction, err := UnmarshalTransaction(file)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, curTransaction)
	}

	return transactions, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
