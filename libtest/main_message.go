package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arstevens/go-libhive-core/message"
	"github.com/arstevens/go-libhive-core/security"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

func main() {
	sk, err := security.RetrieveLocalPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	hMap := make(map[string]interface{})
	hMap[message.TypeField] = message.TransactionType
	header := message.NewMessageHeader(hMap)

	f, err := os.Open("/home/aleksandr/ss")
	if err != nil {
		log.Fatal(err)
	}

	msg, err := message.NewMessage(header, *sk, f)
	if err != nil {
		log.Fatal(err)
	}
	hMap[message.SignField] = "A6XVLywEiBn8bWHl8UfyWtIbhWdTpdohlCOnBEB9AsoHC26675VO/FHTtFpMxJJT85LeVU2B8Glm87r3168qcQLWTORhIRlpxZo0UhoaKBvmgTy/RFC4bJ48VYOXtahaJa/bogbaTNpXqcDfPmtGcObijwbKgJtZ6abXE5jhk8/6wmz/gqzFcGS7u60wSKbdLluD7biaBI7UDH8GljvCYv5HSttysWqNZTHb7XDCMzHniTe1a8VhUCGvYNa+NpT3l0UY99Vx6rYnK+zlB71yiaV7wNFTOFbb1COdqRwYt6Vu4iPNjPv/w9KtWChxd+yDwNgTPoO4MforoJwpvsk8Tw=="
	msg.SetHeader(message.NewMessageHeader(hMap))

	pub := sk.GetPublic()
	x := pub.(*crypto.RsaPublicKey)
	fmt.Println(msg.Verify(*x))

}
