package main

import (
	"fmt"
	"io/ioutil"
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

	data := ""
	buf := make([]byte, 8)
	n, err := msg.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	data += string(buf[:n])
	for n > 0 {
		n, err = msg.Read(buf)
		data += string(buf[:n])
	}

	ioutil.WriteFile("/home/aleksandr/testfile", []byte(data), os.ModePerm)
	f2, _ := os.Open("/home/aleksandr/testfile")
	nmsg, err := message.ReadMessage(f2)
	if err != nil {
		log.Fatal(err)
	}

	n, err = nmsg.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(buf[:n]))
	for n > 0 {
		n, err = nmsg.Read(buf)
		fmt.Print(string(buf[:n]))
	}

}
