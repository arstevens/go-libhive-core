package stream

import (
	"encoding/json"
	"io"

	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p-core/protocol"
	maddr "github.com/multiformats/go-multiaddr"
)

type Stream struct {
	resp io.Close
	dec  *json.Decoder
}

// Run in seperate go routine
func NewStreamHandler(sh *ipfsapi.Shell, callback func(s Stream)) {
}

func NewStream(sh *ipfsapi.Shell) {

}

func streamListen(sh *ipfsapi.Shell, proto protocol.ID, addr maddr.Multiaddr) {

}
