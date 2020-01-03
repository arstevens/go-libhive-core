package main

import (
	"log"

	security "github.com/arstevens/go-libhive-core/security"
)

func main() {
	_, err := security.RetrieveLocalPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
}
