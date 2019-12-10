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
)

type MessageHeader struct {
	header map[string]interface{}
}

// Constructor
func NewMessageHeader(h map[string]interface{}) *MessageHeader {
	msg := MessageHeader{header: h}
	return &msg
}

func NewBufferedMessageHeader(in io.Reader) *MessageHeader {
	bReader := bufio.NewReader(in)
	rawData, err := bReader.ReadSlice(byte(0x03))
	if err != nil {
		fmt.Println("Could not read message from io.Reader object")
		fmt.Println(err.Error())
		return nil
	}

	msg := MessageHeader{header: nil}
	err = msg.Unmarshal(rawData)
	if err != nil {
		fmt.Println("Could not unmarshal bytes read into Message")
		fmt.Println(err.Error())
		return nil
	}
	return &msg
}

// Accessors
func (m *MessageHeader) Type() int {
	mType, ok := (m.header[TypeField]).(int)
	if !ok {
		fmt.Println("Could not assert TypeField to Int in Message")
		return -1
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

func (m *MessageHeader) DataLen() int64 {
	dLen, ok := (m.header[DataLenField]).(int64)
	if !ok {
		fmt.Println("Could not assert DataLenField to Int64 in Message")
		return -1
	}
	return dLen
}

// Marshaling
func (m *MessageHeader) Marshal() []byte {
	jsonString, err := json.Marshal(m.header)
	if err != nil {
		fmt.Println("Could not marshal Header in Message")
		fmt.Println(err.Error())
		return []byte{}
	}
	return []byte(jsonString)
}

func (m *MessageHeader) Unmarshal(raw []byte) error {
	m.header = make(map[string]interface{})
	err := json.Unmarshal(raw, m.header)
	if err != nil {
		fmt.Println("Could not unmarshal data in Message")
		return err
	}
	return nil
}

// Prep for Sending
func PackageBytes(msg []byte) []byte {
	return append(msg, 0x03)
}
