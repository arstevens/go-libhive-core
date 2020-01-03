package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

var MinRsaKeyBits = 2048
var ErrRsaKeyTooSmall error

type RsaPrivateKey struct {
	sk rsa.PrivateKey
}

// Sign returns a signature of the input data
func (sk *RsaPrivateKey) Sign(message []byte) ([]byte, error) {
	hashed := sha256.Sum256(message)
	return rsa.SignPKCS1v15(rand.Reader, &sk.sk, crypto.SHA256, hashed[:])
}

// UnmarshalRsaPrivateKey returns a private key from the input x509 bytes
func UnmarshalRsaPrivateKey(b []byte) (*RsaPrivateKey, error) {
	sk, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	if sk.N.BitLen() < MinRsaKeyBits {
		return nil, ErrRsaKeyTooSmall
	}
	return &RsaPrivateKey{sk: *sk}, nil
}

type RsaPublicKey struct {
	k rsa.PublicKey
}

// Verify compares a signature against input data
func (pk *RsaPublicKey) Verify(data, sig []byte) (bool, error) {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(&pk.k, crypto.SHA256, hashed[:], sig)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UnmarshalRsaPublicKey returns a public key from the input x509 bytes
func UnmarshalRsaPublicKey(b []byte) (*RsaPublicKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		log.Fatal("Could not get pem block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pk, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not actually an rsa public key")
	}
	if pk.N.BitLen() < MinRsaKeyBits {
		return nil, ErrRsaKeyTooSmall
	}
	return &RsaPublicKey{k: *pk}, nil
}

func RetrievePublicKey(sh *ipfsapi.Shell, peerID string) (*RsaPublicKey, error) {
	idOut, err := sh.ID(peerID)
	if err != nil {
		fmt.Println("here1")
		return nil, err
	}
	rawPubKey := []byte("-----BEGIN PUBLIC KEY----- \n" + idOut.PublicKey + "\n-----END PUBLIC KEY-----")

	pubKey, err := UnmarshalRsaPublicKey(rawPubKey)
	fmt.Println("here3")
	return pubKey, err
}
