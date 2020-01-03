package contract

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/arstevens/go-libhive-core/message"
)

const (
	ContractIdField      = "id"
	ContractServiceField = "srv"
	ContractTransField   = "tran"
	ContractEarnField    = "earn"
)

type Contract struct {
	id             string
	service        int
	tokenTransfers map[string]float64
	tokenEarnings  map[string]float64
}

// Constructors
func NewContract(ident string, serv int, tkT map[string]float64, tkE map[string]float64) *Contract {
	c := Contract{id: ident, service: serv, tokenTransfers: tkT, tokenEarnings: tkE}
	return &c
}

func NewBufferedContract(in io.Reader) *Contract {
	bReader := bufio.NewReader(in)
	rawData, err := bReader.ReadSlice(byte(message.EndOfHeader))
	if err != nil {
		fmt.Println("Could not read contract from io.Reader object")
		fmt.Println(err.Error())
		return nil
	}

	c := Contract{}
	err = c.Unmarshal(rawData)
	if err != nil {
		return nil
	}
	return &c
}

// Marshaling
func (c *Contract) Marshal() []byte {
	consolidated := make(map[string]interface{})
	consolidated[ContractIdField] = c.id
	consolidated[ContractServiceField] = c.service
	consolidated[ContractTransField] = c.tokenTransfers
	consolidated[ContractEarnField] = c.tokenEarnings

	marsh, err := json.Marshal(consolidated)
	if err != nil {
		fmt.Println("Could not marshal contract")
		fmt.Println(err.Error())
		return []byte{}
	}
	return []byte(marsh)
}

func (c *Contract) Unmarshal(raw []byte) error {
	consolidated := make(map[string]interface{})
	err := json.Unmarshal(raw, consolidated)
	if err != nil {
		fmt.Println("Could not unmarshal contract")
		fmt.Println(err.Error())
		return err
	}
	c.id = consolidated[ContractIdField].(string)
	c.service = consolidated[ContractServiceField].(int)
	c.tokenTransfers = consolidated[ContractTransField].(map[string]float64)
	c.tokenEarnings = consolidated[ContractEarnField].(map[string]float64)
	return nil
}
