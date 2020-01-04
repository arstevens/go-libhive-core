package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"reflect"
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
	CapsuleField = "cap"
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

func ReadMessageHeader(in *bufio.Reader) (*MessageHeader, error) {
	rawData, err := in.ReadBytes(byte(EndOfHeader))
	fmt.Println(string(rawData))
	rawData = rawData[:len(rawData)-1]
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

func (m *MessageHeader) IsCapsule() bool {
	mType, ok := m.header[CapsuleField].(bool)
	if !ok {
		fmt.Println("Could not assert CapsuleField to bool in message")
		return false
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
	val := reflect.ValueOf(m.header[DataLenField])
	iVal := val.Convert(reflect.TypeOf(0))

	return int(iVal.Int())
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
