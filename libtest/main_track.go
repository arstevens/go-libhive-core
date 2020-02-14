package main

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/arstevens/go-libhive-core/track"
)

func main() {
	parties := make([]track.Party, 10)
	fmt.Println("Parties\n-------")
	for i := 0; i < 10; i++ {
		party, err := createRandomParty()
		if err != nil {
			panic(err)
		}
		fmt.Println(party.Id())
		parties[i] = party
	}

	fmt.Println("Transactions\n-------------")
	for i := 0; i < 10; i += 2 {
		tid := generateRandomString()

		exchanges := make(map[string]float64)
		val := mrand.Float64()
		exchanges[parties[i].Id()] = val
		exchanges[parties[i+1].Id()] = -val
		trans := track.NewTransaction(tid, exchanges, time.Now())
		fmt.Println(exchanges)
		err := parties[i].AddTransaction(*trans)
		if err != nil {
			panic(err)
		}
	}

}

func createRandomParty() (track.Party, error) {
	pid := generateRandomString()
	party := track.NewParty(pid, "/home/aleksandr/fsLoc", 0.0)
	return *party, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}
