package security

import (
	"encoding/base64"

	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func RetrievePublicKey(sh *ipfsapi.Shell, peerID string) (crypto.PubKey, error) {
	idOut, err := sh.ID(peerID)
	if err != nil {
		return nil, err
	}
	rawPubKey := idOut.PublicKey
	protoPubKey, err := base64.StdEncoding.DecodeString(rawPubKey)
	if err != nil {
		return nil, err
	}

	pubKey, err := crypto.UnmarshalPublicKey(protoPubKey)
	return pubKey, err
}
