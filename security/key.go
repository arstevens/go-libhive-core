package security

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	crypto "github.com/libp2p/go-libp2p-core/crypto"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

func RetrievePublicKey(sh *ipfsapi.Shell, peerID string) (*crypto.RsaPublicKey, error) {
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
	return pubKey.(*crypto.RsaPublicKey), err
}

func RetrieveLocalPrivateKey() (*crypto.RsaPrivateKey, error) {
	var result map[string]interface{}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(home + "/.ipfs/config")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	identity := result["Identity"].(map[string]interface{})
	sk := identity["PrivateKey"].(string)

	protoSK, err := base64.StdEncoding.DecodeString(sk)
	if err != nil {
		return nil, err
	}

	privKey, err := crypto.UnmarshalPrivateKey(protoSK)
	return privKey.(*crypto.RsaPrivateKey), err
}
