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
	defer f.Close()

	msg, err := message.NewMessage(header, sk, f)
	if err != nil {
		log.Fatal(err)
	}

	h2 := make(map[string]interface{})
	h2[message.TypeField] = message.QueryType

	msg2, err := message.NewMessage(message.NewMessageHeader(h2), sk, msg)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 8)
	n, err := msg2.Read(buf)
	fmt.Print(string(buf[:n]))
	for n > 0 {
		n, err = msg2.Read(buf)
		fmt.Print(string(buf[:n]))
	}

	fmt.Println("\n---------------------- DECAPSULATE RUN ----------------------")
	msgs, err := message.Decapsulate(msg2)
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range msgs {
		fmt.Println(m.Header().DataLen())
	}
}
