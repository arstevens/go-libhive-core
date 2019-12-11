package message

import (
	"fmt"
	"io"
	"io/ioutil"
	"bufio"
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
	bodyLen := h.DataLen()
	tRead := int64(0)
	rHead := h.Marshal()
	file, _ := ioutil.TempFile(os.TempDir(), "msg")
	buf := make([]byte, 8192)
	var n int
	if bodyLen < len(buf) {
		n, _ = r.Read(buf[:bodyLen])
	} else {
		n, _ = r.Read(buf)
	}
	tRead += n
	for n > 0 && tRead < bodyLen {
		_, err := file.Write(buf[:n])
		if err != nil {
			fmt.Println("Could not write to temp file to create message")
			return nil
		}
		if (bodyLen - tRead) < len(buf) {
			n, _ = r.Read(buf[:bodyLen - tRead])
		} else {
			n, _ = r.Read(buf)
		}
	}
	file.Seek(0, io.SeekStart)

	m := Message{bytePtr: 0, header: h, rawHeader: rHead, body: file, bodyBufFile: file.Name()}
	return &m
}

func (m *Message) Header() *MessageHeader {
	return m.header
}

func (m *Message) SetHeader(h *MessageHeader) {
	m.header = h
}

func (m *Message) Body() *os.File {
	return m.body
}

func (m *Message) SetBody(r io.Reader) {
	m = NewMessage(m.header, r)
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

func (m *Message) ReadUntil(b byte) ([]byte, error) {
	r := bufio.NewReader(m.body)
	return r.ReadSlice(b)
}

func (m *Message) Clean() error {
	return m.body.Close()
}
