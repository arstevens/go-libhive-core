package main

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/arstevens/go-libhive-core/track"
)

func main() {
	/*
		  // Party creation and transaction Marshalling test
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
	*/

	// Transaction Unmarshalling test
	/*
		f, err := os.Open("/home/aleksandr/fsLoc/1581714437830101486/1581714437830603015")
		if err != nil {
			panic(err)
		}
		t, err := track.UnmarshalTransaction(f)
		if err != nil {
			panic(err)
		}
		fmt.Println(t.Id())
		fmt.Println(t.Parties())
		fmt.Println(t.Time())
	*/

	/*
		// Party Transaction Sum test
		p := track.NewParty("1", "/home/aleksandr/fsLoc", 0.0)
		exchanges := make(map[string]float64)
		exchanges["1"] = 0.1
		exchanges["2"] = -0.1
		t1 := track.NewTransaction("a", exchanges, time.Now())
		p.AddTransaction(*t1)
		delete(exchanges, "2")
		exchanges["3"] = -0.1
		t2 := track.NewTransaction("b", exchanges, time.Now())
		p.AddTransaction(*t2)

		sum, err := p.SumTransactions()
		if err != nil {
			panic(err)
		}
		fmt.Println(sum)
	*/

	//DAG test
	eg, err := track.NewExchangeGraph("/home/aleksandr/fsLoc")
	if err != nil {
		panic(err)
	}

	p1 := eg.GetParty("1")
	p2 := eg.GetParty("2")
	p3 := eg.GetParty("3")
	fmt.Println(p1.Id())
	fmt.Println(p2.Id())
	fmt.Println(p3.Id())

	err = eg.Compress("/home/aleksandr/fsLoc/history.json")
	if err != nil {
		panic(err)
	}
	err = eg.DeleteTransactions()
	if err != nil {
		panic(err)
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
