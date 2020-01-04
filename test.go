package main

import (
	"bufio"
	"log"
	"os"
)

func hasEOT(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] == 0x04 {
			return true
		}
	}
	return false
}

func main() {
	file, err := os.Open("/home/aleksandr/testfile")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	/*
		buf := make([]byte, 16)
		n, _ := file.Read(buf)
		fmt.Println(hasEOT(buf[:n]))
		for n > 0 {
			n, _ = file.Read(buf)
			fmt.Println(hasEOT(buf[:n]))
		}
	*/

	reader := bufio.NewReader(file)
	_, err = reader.ReadBytes(0x04)
	if err != nil {
		log.Fatal(err)
	}

}
