package node

import (
	"github.com/arstevens/go-libhive-core/stream"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

type Node struct {
	conn *stream.Stream
	pk   *crypto.RsaPublicKey
}

func (n *Node) Key() *crypto.RsaPublicKey {
	return n.pk
}

func (n *Node) Conn() *stream.Stream {
	return n.conn
}
