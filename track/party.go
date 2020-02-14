package track

import (
	"log"
	"os"
	"path/filepath"
)

// Is there a way to have nodes only keep info on selective parties &
// still enable synchronous consensus?
type Party struct {
	id         string
	fsLocation string
	history    float64
	parser     TransactionParser
}

func (p *Party) NewParty(id string, loc string, hist float64) *Party {
	return &Party{id: id, fsLocation: loc, history: hist}
}

func (p *Party) Id() string {
	return p.id
}

func (p *Party) AddTransaction(t Transaction) error {
	// Add file to the folder for this parties transactions
	transactionFile, err := os.Open(p.fsLocation + "/" + t.Id())
	if err != nil {
		return err
	}
	defer transactionFile.Close()

	serial := t.Marshal()
	_, err = transactionFile.Write(serial)

	// Link this transaction to the other members of this transaction
	parties := t.Parties()
	parentDir := filepath.Dir(p.fsLocation)
	recordedParties := readDirectory(parentDir)

	for _, party := range parties {
		partyFound := false
		for i := 0; i < len(recordedParties) && !partyFound; i++ {
			if party == recordedParties[i] {
				newRoot := parentDir + "/" + party
				err = os.Symlink(transactionFile.Name(), newRoot+"/"+t.Id())
				if err != nil {
					log.Fatal(err)
				}

				partyFound = true
			}
		}
		if !partyFound {
			newRoot := parentDir + "/" + party
			err = os.Mkdir(newRoot, os.ModeDir)
			if err != nil {
				log.Fatal(err)
			}

			err = os.Symlink(transactionFile.Name(), newRoot+"/"+t.Id())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return err
}

func (p *Party) SumTransactions() (float64, error) {
	transactionPaths := readDirectory(p.fsLocation)
	transcations, err := p.parser(transactionPaths)
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

// Should be moved to go-libhive and refactored
/*
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
*/
