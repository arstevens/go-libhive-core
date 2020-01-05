package track

import (
	"log"
	"os"
	"path/filepath"
)

type Party struct {
	id         string
	fsLocation string
}

func (p Party) AddTransaction(t Transaction) error {
	// Add file to the folder for this parties transactions
	transactionFile, err := os.Open(p.fsLocation + "/" + t.transactionId)
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
				err = os.Symlink(transactionFile.Name(), newRoot+"/"+t.transactionId)
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

			err = os.Symlink(transactionFile.Name(), newRoot+"/"+t.transactionId)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return err
}

func (p Party) SumTransactions() (float64, error) {
	transactionPaths := readDirectory(p.fsLocation)
	transcations, err := parseTransactions(transactionPaths)
	if err != nil {
		return -1.0, err
	}

	sum := 0.0
	for _, curTransaction := range transcations {
		sum += curTransaction.GetAmountExchanged(p.id)
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

		curTransaction, err := ParseTransaction(file)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, curTransaction)
	}

	return transactions, nil
}
