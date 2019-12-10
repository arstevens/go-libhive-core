package message

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Message struct {
	header      *MessageHeader
	rawHeader   []byte
	bodyBufFile string
	body        *os.File
	bytePtr     int
}

func NewMessage(h *MessageHeader, r io.Reader) *Message {
	rHead := h.Marshal()
	file, _ := ioutil.TempFile(os.TempDir(), "msg")
	buf := make([]byte, 8192)
	n, _ := r.Read(buf)
	for n > 0 {
		_, err := file.Write(buf[:n])
		if err != nil {
			fmt.Println("Could not write to temp file to create message")
			return nil
		}
		n, _ = r.Read(buf)
	}
	file.Seek(0, io.SeekStart)

	m := Message{bytePtr: 0, header: h, rawHeader: rHead, body: file, bodyBufFile: file.Name()}
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

func (m *Message) Clean() error {
	return m.body.Close()
}
