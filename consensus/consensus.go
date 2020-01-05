package consensus

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

func accumulateSubnetResponse(subnet *Group, msg *message.Message) ([]*message.Message, error) {
	// Start communication loop
	_, entryConn := subnet.EntryNode()
	err := entryConn.WriteReader(msg)
	defer entryConn.Close()

	if err != nil {
		return nil, err
	}

	// End communication loop
	_, exitConn := subnet.ExitNode()
	respMsg, err := message.ReadMessage(exitConn)
	if err != nil {
		return nil, err
	}

	// Parse raw bytes
	return message.Decapsulate(respMsg)
}

// @return: % of verified subnet consensus, % of subnet verified, error
func BinaryConsensus(subnet *Group, sk *crypto.RsaPrivateKey, value []byte) (float32, float32, error) {
	// Create Message
	hMap := make(map[string]interface{})
	hMap[message.TypeField] = ConsensusType
	hMap[message.DataLenField] = len(value)
	header := message.NewMessageHeader(hMap)
	valueReader := bytes.NewReader(value)
	msg, err := message.NewMessage(header, sk, valueReader)

	// Grab public keys from subnet
	// TODO: May want to run this in seperate goroutine while waiting for ECL
	nodes := subnet.SortedKeys()
	keys := make([]*crypto.RsaPublicKey, len(nodes))
	for i, k := range nodes {
		pubKey, err := security.RetrievePublicKey(subnet.GetShell(), k)
		if err != nil {
			return 0.0, 0.0, err
		}
		keys[i] = pubKey
	}

	// Retrieve messages from the network
	layers, err := accumulateSubnetResponse(subnet, msg)

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
