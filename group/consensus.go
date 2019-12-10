package group

import (
	"bytes"

	"github.com/arstevens/go-libhive-core/message"
)

const (
	ConsensusType = "consensus"
)

// @return: % of subnet consensus
func Consensus(subnet *Group, value []byte) (float32, error) {
	// Create Message
	hMap := make(map[string]interface{})
	hMap[message.TypeField] = ConsensusType
	hMap[message.DataLenField] = len(value)
	header := message.NewMessageHeader(hMap)
	valueReader := bytes.NewReader(value)
	msg := message.NewMessage(header, valueReader)

	_, entryConn := subnet.EntryNode()
	err := conn.WriteReader(msg)
	defer entryConn.Close()

	_, exitConn := subnet.ExitNode()
	ch := exitConn.ChanReader()
	// must buffer until reach EndOfTransmission character
	resp := <-ch

	// Unmarshal into Message
	// Ensure each signature is valid in value custody chain
	// compute consensus score
}
