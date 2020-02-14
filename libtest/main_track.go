package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	mrand "math/rand"
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
		tid, err := generateRandomString(int(time.Now().UnixNano()))
		if err != nil {
			panic(err)
		}

		exchanges := make(map[string]float64)
		val := mrand.Float64()
		exchanges[parites[i].Id()] = val
		exchanges[parties[i+1].Id()] = -val
		trans := track.NewTransaction(tid, exchanges, time.Now())
		fmt.Println(exchanges)
		err = parties[i].AddTransaction(trans)
		if err != nil {
			panic(err)
		}
	}

}

func createRandomParty() (track.Party, error) {
	pid, err := generateRandomString(int(time.Now().UnixNano()))
	if err != nil {
		return track.Party{}, err
	}

	party := track.NewParty(pid, "/home/aleksandr/fsLoc", 0.0)
	return party, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
