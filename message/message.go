package message

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

const (
	MsgBufSize = 8192
)

type Message struct {
	header    *MessageHeader
	rawHeader []byte
	body      *os.File
	readPtr   int
}

func NewMessage(h *MessageHeader, sk *crypto.RsaPrivateKey, r io.Reader) (*Message, error) {
	// Create buffer file used for message resets
	bodyFile, err := ioutil.TempFile(os.TempDir(), "msgbuf")
	if err != nil {
		return nil, err
	}

	// Write data to buffer file
	buf := make([]byte, MsgBufSize)
	tRead := 0
	nIn, _ := r.Read(buf)

	_, err = bodyFile.Write(buf[:nIn])
	if err != nil {
		return nil, err
	}

	tRead += nIn
	for nIn > 0 {
		nIn, _ = r.Read(buf)
		_, err = bodyFile.Write(buf[:nIn])
		if err != nil {
			return nil, err
		}
		tRead += nIn
	}
	bodyFile.Seek(0, io.SeekStart)

	// Generate Sign of message for verification by other nodes
	hash, err := hashFile(bodyFile)
	if err != nil {
		return nil, err
	}
	bodyFile.Seek(0, io.SeekStart)
	sign, err := sk.Sign(hash)
	if err != nil {
		return nil, err
	}
	h.header[SignField] = base64.StdEncoding.EncodeToString(sign)

	// Prepare header for message
	h.header[DataLenField] = tRead
	rHeader := h.Marshal()

	m := Message{readPtr: 0, header: h, rawHeader: rHeader, body: bodyFile}
	return &m, nil
}

func ReadMessage(in io.Reader) (*Message, error) {
	// Read Message header first
	head, err := ReadMessageHeader(in)
	if err != nil {
		return nil, err
	}

	// Read remaining data
	bodyFile, err := ioutil.TempFile(os.TempDir(), "msgbuf")
	if err != nil {
		return nil, err
	}

	bodyLen := head.DataLen()
	buf := make([]byte, MsgBufSize)
	tRead := 0
	nIn := 0

	for tRead < bodyLen {
		// If buffer is bigger than needed shrink it
		bytesLeft := bodyLen - tRead
		if bytesLeft > MsgBufSize {
			buf = buf[:bytesLeft]
		}

		nIn, _ = in.Read(buf)
		_, err = bodyFile.Write(buf[:nIn])
		if err != nil {
			return nil, err
		}
		tRead += nIn
		fmt.Println(tRead)
	}
	bodyFile.Seek(0, io.SeekStart)

	m := Message{header: head, rawHeader: head.Marshal(), body: bodyFile, readPtr: 0}
	return &m, nil
}

func (m *Message) Verify(k *crypto.RsaPublicKey) bool {
	// Calculate a hash value for the body
	m.Reset()
	hash, err := hashFile(m.body)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve sign from header
	sign, err := base64.StdEncoding.DecodeString(m.header.MessageSign())
	if err != nil {
		log.Fatal(err)
	}

	// Verify the message
	verf, _ := k.Verify(hash, sign)
	return verf
}

// Close should be defered after Message is created
func (m *Message) Close() error {
	return m.body.Close()
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

func (m *Message) Reset() error {
	_, err := m.body.Seek(0, io.SeekStart)
	m.readPtr = 0
	return err
}

func (m *Message) Read(b []byte) (int, error) {
	// Decide whether to read from header or body
	if m.readPtr < len(m.rawHeader) {
		// Read strictly from header
		if (m.readPtr + len(b)) < len(m.rawHeader) {
			copy(b, m.rawHeader[m.readPtr:m.readPtr+len(b)])
			m.readPtr += len(b)
			return len(b), nil
		} else { // Read partially from header and partially from body
			copy(b, m.rawHeader[m.readPtr:len(m.rawHeader)])
			bufIdx := len(m.rawHeader) - m.readPtr
			n, err := m.body.Read(b[bufIdx:])
			if err != nil {
				return n, err
			}
			m.readPtr += bufIdx + n
			return n + bufIdx, nil
		}
	} else { // Read from body
		n, err := m.body.Read(b)
		m.readPtr += n
		return n, err
	}
}

func hashFile(f *os.File) ([]byte, error) {
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
