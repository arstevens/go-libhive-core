package security

import (
  ipfsapi "github.com/ipfs/go-ipfs-api"
  "github.com/libp2p/go-libp2p-core/crypto"
  "encoding/base64"
)

func RetrievePublicKey(sh *ipfsapi.Shell, peerID string) (crypto.RsaPublicKey, error) {
  idOut, err := sh.ID(peerID)
  if err != nil {
    return nil, err
  }
  rawPubKey := idOut.PublicKey
  protoPubKey, err := base64.DecodeString(rawPublicKey)
  if err != nil {
    return nil, err
  }

  pubKey, err := crypto.UnmarshalPublicKey(protoPubKey)
  return (*pubKey).(crypto.RsaPublicKey), err
}
