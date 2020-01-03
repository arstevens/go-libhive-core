package main

import (
	"log"

	security "github.com/arstevens/go-libhive-core/security"
	ipfsapi "github.com/ipfs/go-ipfs-api"
)

func main() {
	sh := ipfsapi.NewLocalShell()
	_, err := security.RetrievePublicKey(sh, "QmadRNcr9cdaqaDFCmj7VfSTzUN6pwZFCs7952bRTnE5qF")
	if err != nil {
		log.Fatal(err)
	}
}
