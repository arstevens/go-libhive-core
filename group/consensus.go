package group

import (
	"bytes"

	"github.com/arstevens/go-libhive-core/message"
	"github.com/arstevens/go-libhive-core/security"
	"github.com/libp2p/go-libp2p-core/crypto"
)

const (
	ConsensusType   = "consensus"
	PropogationType = "propogation"
)

// @return: % of verified subnet consensus, % of subnet verified, error
func BasicConsensus(subnet *Group, value []byte) (float32, float32, error) {
	// Create Message
	hMap := make(map[string]interface{})
	hMap[message.TypeField] = ConsensusType
	hMap[message.DataLenField] = len(value)
	header := message.NewMessageHeader(hMap)
	valueReader := bytes.NewReader(value)
	msg := message.NewMessage(header, valueReader)

	// Start communication loop
	_, entryConn := subnet.EntryNode()
	err := conn.WriteReader(msg)
	defer entryConn.Close()

	// End communication loop
	_, exitConn := subnet.ExitNode()
	rHeader := NewBufferedMessageHeader(exitConn)
	respMsg := NewMessage(rHeader, exitConn)

	// Grab public keys from subnet
	// TODO: May want to run this in seperate goroutine while waiting for ECL
	nodes := subnet.SortedKeys()
	keys := make([]crypto.RsaPublicKey, len(nodes))
	for i, k := range nodes {
		pubKey, err := security.RetrievePublicKey(g.GetShell(), k)
		if err != nil {
			return 0.0, 0.0, err
		}
		keys[i] = pubKey
	}

	// Parse raw bytes
	layers, err := message.Decapsulate(respMsg)
	signedValues := make([]message.SignedValue, len(nodes))
	for _, capsule := range layers {
		signedValues[i] = capsule.SignedValue()
	}

	// Calculate consensus scores
	nVerified := 0
	nConsensus := 0
	for i, pKey := range keys {
		sVal := signedValues[i]
		verified, err := pKey.Verify(sVal.value, sVal.sign)
		if err != nil {
			return nVerified, nConsensus, err
		}

		if verified {
			nVerified += 1
			if bytes.Equal(value, sVal.value) {
				nConsensus += 1
			}
		}
	}

	// Propogate recieved (value, sign) tuples to subnet for individual consensus
	respMsg.Reset()
	pmHeaderMap := make(map[string]interface{})
	pmHeaderMap[message.TypeField] = PropogationType
	pmHeaderMap[message.DataLenField] = rHeader.DataLen()
	respMsg.setHeader(message.NewMessageHeader(pmHeaderMap))

	// add origin/chain of custody to reduce redundant transmissions
	// Use message encapsulation for Chain of Custody?
	err = entryConn.WriteReader(respMsg)
	if err != nil {
		return 0.0, 0.0, err
	}
	eHeader := NewBufferedMessageHeader(exitConn)

	pVerified := float32(nVerified) / float32(len(nodes))
	pConsensus := float32(nConsensus) / float32(len(nodes))
	return pConsensus, pVerified, nil
}
