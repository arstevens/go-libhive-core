package message

import (
	"fmt"
	"os"
)

type Message struct {
	header      MessageHeader
	rawHeader   []byte
	bodyBufFile string
	body        *os.File
	bytePtr     int
}

func NewMessage(h MessageHeader, tmpf string) *Message {
	rHead := h.Marshal()
	file, err := os.Open(tmpf)
	if err != nil {
		fmt.Println("Could not read temp file for new message")
		fmt.Println(err.Error())
		return nil
	}
	m := Message{bytePtr: 0, header: h, rawHeader: rHead, body: file, bodyBufFile: tmpf}
	return &m
}

func (m *Message) Reset() error {
	err := m.body.Close()
	if err != nil {
		return err
	}

	m.body, err = os.Open(m.bodyBufFile)
	return err
}

func (m *Message) Read(b []byte) (int, error) {
	if m.bytePtr < len(m.rawHeader) {
		if (m.bytePtr + len(b)) < len(m.rawHeader) {
			copy(b, m.rawHeader[m.bytePtr:m.bytePtr+len(b)])
			m.bytePtr += len(b)
			return len(b), nil
		} else {
			copy(b, m.rawHeader[m.bytePtr:len(m.rawHeader)])
			bufIdx := len(m.rawHeader) - m.bytePtr
			n, err := m.body.Read(b[bufIdx:])
			if err != nil {
				return n, err
			}
			m.bytePtr += bufIdx + n
			return n, nil
		}
	} else {
		n, err := m.body.Read(b)
		m.bytePtr += n
		return n, err
	}
}
