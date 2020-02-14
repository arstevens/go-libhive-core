package track

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"
)

type Transaction struct {
	transactionId string
	exchanges     map[string]float64
	gmtTimestamp  time.Time // should be GMT
}

func NewTransaction(tid string, exchanges map[string]float64, tstamp time.Time) *Transaction {
	return &Transaction{transactionId: tid, exchanges: exchanges, gmtTimestamp: tstamp}
}

func (t Transaction) Id() string {
	return t.transactionId
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

func (t Transaction) Marshal() []byte {
	serial, err := json.Marshal(t)
	fmt.Println("Serial: " + string(serial))
	if err != nil {
		log.Fatal(err)
	}

	return serial

	/*
		serial := t.transactionId + ","
		exchangeSerial := ""
		for party, value := range t.exchanges {
			fStr := strconv.FormatFloat(value, 'E', -1, 64)
			exchangeSerial += party + ":" + fStr + ","

		}
	*/
}

func UnmarshalTransaction(r io.Reader) (*Transaction, error) {
	var newTransaction Transaction
	serial, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(serial, &newTransaction)
	return &newTransaction, err
}
