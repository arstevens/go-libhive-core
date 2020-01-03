package main

import (
	"log"

	security "github.com/arstevens/go-libhive-core/security"
	ipfsapi "github.com/ipfs/go-ipfs-api"
)

func main() {
	sh := ipfsapi.NewLocalShell()
	_, err := security.RetrieveLocalPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
}
