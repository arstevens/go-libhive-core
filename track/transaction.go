package track

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

type Transaction struct {
	transactionId string
	exchanges     map[string]float64
}

// Returns everyone involved in the transaction
func (t Transaction) Parties() []string {
	parties := make([]string, len(t.exchanges))
	i := 0
	for k, _ := range t.exchanges {
		parties[i] = k
		i++
	}

	return parties
}

func (t Transaction) GetAmountExchanged(id string) float64 {
	return t.exchanges[id]
}

// Serialize a transaction for filesystem storage
func (t Transaction) Marshal() []byte {
	serial, err := json.Marshal(t.exchanges)
	if err != nil {
		log.Fatal(err)
	}

	return serial
}

func ParseTransaction(r io.Reader) (*Transaction, error) {
	var newTransaction *Transaction
	serial, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(serial, &newTransaction.exchanges)
	return newTransaction, err
}
