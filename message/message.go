package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

const (
	TransactionType = iota
	QueryType       = iota
)

const (
	TypeField    = "trans"
	ConvIdField  = "convid"
	MsgIdField   = "msgid"
	DataLenField = "dlen"
)

type Message struct {
	header map[string]interface{}
}

// Constructor
func NewMessage(h map[string]interface{}) *Message {
	msg := Message{header: h}
	return &msg
}

func NewBufferedMessage(in io.Reader) *Message {
	bReader := bufio.NewReader(in)
	rawData, err := bReader.ReadSlice(byte(0x03))
	if err != nil {
		fmt.Println("Could not read message from io.Reader object")
		fmt.Println(err.Error())
		return nil
	}

	msg := Message{header: nil}
	err = msg.Unmarshal(rawData)
	if err != nil {
		fmt.Println("Could not unmarshal bytes read into Message")
		fmt.Println(err.Error())
		return nil
	}
	return &msg
}

// Accessors
func (m *Message) Type() int {
	mType, ok := (m.header[TypeField]).(int)
	if !ok {
		fmt.Println("Could not assert TypeField to Int in Message")
		return -1
	}

	return mType
}

func (m *Message) ConvId() string {
	convid, ok := (m.header[ConvIdField]).(string)
	if !ok {
		fmt.Println("Could not assert ConvIdField to String in Message")
		return ""
	}
	return convid
}

func (m *Message) MsgId() int {
	mId, ok := (m.header[MsgIdField]).(int)
	if !ok {
		fmt.Println("Could not assert MsgIdField to Int in Message")
		return -1
	}
	return mId
}

func (m *Message) DataLen() int64 {
	dLen, ok := (m.header[DataLenField]).(int64)
	if !ok {
		fmt.Println("Could not assert DataLenField to Int64 in Message")
		return -1
	}
	return dLen
}

// Marshalling
func (m *Message) Marshal() []byte {
	jsonString, err := json.Marshal(m.header)
	if err != nil {
		fmt.Println("Could not marshal Header in Message")
		fmt.Println(err.Error())
		return []byte{}
	}
	return []byte(jsonString)
}

func (m *Message) Unmarshal(raw []byte) error {
	m.header = make(map[string]interface{})
	err := json.Unmarshal(raw, m.header)
	if err != nil {
		fmt.Println("Could not unmarshal data in Message")
		return err
	}
	return nil
}
