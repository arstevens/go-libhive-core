package track

import "io"

// need to make transaction an interface to allow different type of
// datapoints to be stored(aka WDS information)
type Transaction interface {
	Id() string
	Parties() []string
	Marshal() []byte
	Unmarshal(io.Reader) error
	GetAmountExchanged(string) float64
}

type TransactionParser func(paths []string) ([]*Transaction, error)

/*
type Transaction struct {
	transactionId string
	exchanges     map[string]float64
	gmtTimestamp  time.Time // should be GMT
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
	if err != nil {
		log.Fatal(err)
	}

	return serial
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
*/
