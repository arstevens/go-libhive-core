package main

import (
	"fmt"
	"log"

	security "github.com/arstevens/go-libhive-core/security"
)

func main() {
	sk, err := security.RetrieveLocalPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sk)
}
