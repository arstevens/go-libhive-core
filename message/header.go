package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

const (
	TransactionType = "transaction"
	QueryType       = "query"
)

const (
	TypeField    = "trans"
	ConvIdField  = "convid"
	MsgIdField   = "msgid"
	DataLenField = "dlen"
	SignField    = "sign"
)

const (
	EndOfHeader = 0x04
)

type MessageHeader struct {
	header map[string]interface{}
}

// Constructor
func NewMessageHeader(h map[string]interface{}) *MessageHeader {
	msg := MessageHeader{header: h}
	return &msg
}

func ReadMessageHeader(in io.Reader) (*MessageHeader, error) {
	bReader := bufio.NewReader(in)
	rawData, err := bReader.ReadBytes(byte(EndOfHeader))
	rawData = rawData[:len(rawData)-1]
	fmt.Println(rawData)
	if err != nil {
		fmt.Println("Could not read message from io.Reader object")
		return nil, err
	}

	msg := MessageHeader{header: nil}
	err = msg.Unmarshal(rawData)
	if err != nil {
		fmt.Println("Could not unmarshal bytes read into Message")
		fmt.Println(err.Error())
		return nil, err
	}
	return &msg, nil
}

// Accessors
func (m *MessageHeader) Type() string {
	mType, ok := (m.header[TypeField]).(string)
	if !ok {
		fmt.Println("Could not assert TypeField to Int in Message")
		return ""
	}

	return mType
}

func (m *MessageHeader) ConvId() string {
	convid, ok := (m.header[ConvIdField]).(string)
	if !ok {
		fmt.Println("Could not assert ConvIdField to String in Message")
		return ""
	}
	return convid
}

func (m *MessageHeader) MsgId() int {
	mId, ok := (m.header[MsgIdField]).(int)
	if !ok {
		fmt.Println("Could not assert MsgIdField to Int in Message")
		return -1
	}
	return mId
}

func (m *MessageHeader) DataLen() int {
	dLen, ok := m.header[DataLenField].(int)
	fmt.Println(m.header[DataLenField])
	dLen2, ok := m.header[DataLenField].(string)
	fmt.Println(ok)
	fmt.Println(dLen2)
	if !ok {
		fmt.Println("Could not assert DataLenField to Int in Message")
		return -1
	}
	return dLen
}

func (m *MessageHeader) MessageSign() string {
	sVal, ok := m.header[SignField].(string)
	if !ok {
		fmt.Println("Could not assert SignField")
		return ""
	}
	return sVal
}

// Marshaling
func (m *MessageHeader) Marshal() []byte {
	jsonString, err := json.Marshal(m.header)
	if err != nil {
		fmt.Println("Could not marshal Header in Message")
		fmt.Println(err.Error())
		return []byte{}
	}
	return PackageBytes([]byte(jsonString))
}

func (m *MessageHeader) Unmarshal(raw []byte) error {
	m.header = make(map[string]interface{})
	err := json.Unmarshal(raw, &m.header)
	if err != nil {
		fmt.Println("Could not unmarshal data in Message")
		return err
	}
	return nil
}

// Prep for Sending
func PackageBytes(msg []byte) []byte {
	return append(msg, EndOfHeader)
}
