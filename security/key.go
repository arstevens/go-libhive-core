package security

import (
	"encoding/base64"
	"fmt"

	crypto "github.com/libp2p/go-libp2p-core/crypto"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

func RetrievePublicKey(sh *ipfsapi.Shell, peerID string) (*crypto.RsaPublicKey, error) {
	idOut, err := sh.ID(peerID)
	if err != nil {
		fmt.Println("here1")
		return nil, err
	}
	rawPubKey := idOut.PublicKey
	protoPubKey, err := base64.StdEncoding.DecodeString(rawPubKey)
	if err != nil {
		fmt.Println("here2")
		return nil, err
	}

	pubKey, err := crypto.UnmarshalPublicKey(protoPubKey)
	fmt.Println("here3")
	return pubKey.(*crypto.RsaPublicKey), err
}
