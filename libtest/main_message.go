package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arstevens/go-libhive-core/message"
	"github.com/arstevens/go-libhive-core/security"
)

func main() {
	sk, err := security.RetrieveLocalPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	hMap := make(map[string]interface{})
	hMap[message.TypeField] = message.TransactionType
	header := message.NewMessageHeader(hMap)

	f, err := os.Open("/home/aleksandr/index.html")
	if err != nil {
		log.Fatal(err)
	}

	msg, err := message.NewMessage(header, sk, f)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 8)
	n, err := msg.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(buf[:n]))
	for n > 0 {
		n, err = msg.Read(buf)
		fmt.Print(string(buf[:n]))
	}

}
