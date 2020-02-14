package track

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
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

type test struct {
	v1 string
	v2 string
}

func (t Transaction) Marshal() []byte {
	serial := t.transactionId + ","
	for party, value := range t.exchanges {
		fStr := strconv.FormatFloat(value, 'E', -1, 64)
		serial += party + ":" + fStr + ","
	}

	timeMarshal, err := t.gmtTimestamp.MarshalBinary()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	serial += string(timeMarshal)
	return []byte(serial)
}

func UnmarshalTransaction(r io.Reader) (*Transaction, error) {
	var newTransaction Transaction
	serial, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	serialParts := strings.Split(string(serial), ",")

	exchanges := make(map[string]float64)
	for i := 1; i < len(serialParts)-1; i++ {
		exchangeParts := strings.Split(serialParts[i], ":")
		exVal, err := strconv.ParseFloat(exchangeParts[1], 64)
		if err != nil {
			return nil, err
		}
		exchanges[exchangeParts[0]] = exVal
	}

	fmt.Println(serialParts[0])
	fmt.Println(exchanges)

	tTime := time.Time{}
	timeBytes := []byte(serialParts[len(serialParts)-1])
	fmt.Println(timeBytes)
	err = tTime.UnmarshalBinary(timeBytes)
	if err != nil {
		fmt.Println("here")
		return nil, err
	}

	newTransaction.transactionId = serialParts[0]
	newTransaction.exchanges = exchanges
	newTransaction.gmtTimestamp = tTime

	err = json.Unmarshal(serial, &newTransaction)
	return &newTransaction, err
}
