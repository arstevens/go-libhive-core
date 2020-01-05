package node

import (
	"github.com/arstevens/go-libhive-core/stream"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

// Consolodate values associated with a node
type Node struct {
	Id     string
	Stream *stream.Stream
	PubKey *crypto.RsaPublicKey
}

// Slice of Nodes that implements sort.Interface
type Nodes []Node

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Less(i int, j int) bool {
	return ns[i].Id < ns[j].Id
}

func (ns Nodes) Swap(i int, j int) {
	tmp := ns[i]
	ns[i] = ns[j]
	ns[j] = tmp
}
