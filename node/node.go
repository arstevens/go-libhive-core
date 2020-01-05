package node

import (
	"github.com/arstevens/go-libhive-core/stream"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

// Interface to allow comparison between nodes
type Node interface {
	Id() string
}

// Structure to represent local node
type LocalNode struct {
	id      string
	privKey *crypto.RsaPrivateKey
}

func (ln LocalNode) Id() string {
	return ln.id
}

func (ln LocalNode) PrivKey() *crypto.RsaPrivateKey {
	return ln.privKey
}

// Structure to represent nodes controlled by a seperate party
type RemoteNode struct {
	id     string
	stream *stream.Stream
	pubKey *crypto.RsaPublicKey
}

func (rn RemoteNode) Id() string {
	return rn.id
}

func (rn RemoteNode) Stream() *stream.Stream {
	return rn.stream
}

func (rn RemoteNode) PubKey() *crypto.RsaPublicKey {
	return rn.pubKey
}

// Slice of Nodes that implements sort.Interface
type Nodes []Node

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Less(i int, j int) bool {
	return ns[i].Id() < ns[j].Id()
}

func (ns Nodes) Swap(i int, j int) {
	tmp := ns[i]
	ns[i] = ns[j]
	ns[j] = tmp
}

// Returns the nodes with Node Ids on the left and right of you
func (ns Nodes) FindNeighbors(n Node) (*Node, *Node) {
	// Handle edge cases
	if ns.Len() < 1 {
		return nil, nil
	} else if ns.Len() == 1 {
		if ns[0].Id() < n.Id() {
			return &ns[0], nil
		}
		return nil, &ns[0]
	}

	// Find Neighbors
	var left Node
	var right Node
	for _, node := range ns {
		if node.Id() > n.Id() {
			if right.Id() == "" {
				right = node
			} else if node.Id() < right.Id() {
				right = node
			}
		} else if node.Id() < n.Id() {
			if left.Id() == "" {
				left = node
			} else if node.Id() > left.Id() {
				left = node
			}
		}
	}

	return &left, &right
}
