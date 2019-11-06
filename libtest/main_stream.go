package main

import (
	"fmt"
	"os"

	hproto "github.com/arstevens/go-libhive-core/protocol"
	hstream "github.com/arstevens/go-libhive-core/stream"
	ipfsapi "github.com/ipfs/go-ipfs-api"
)

func main() {
	sh := ipfsapi.NewLocalShell()
	proto := hproto.NewProtocolId("htest/1.0")
	hs, err := hstream.NewStream(sh, proto, os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	buf := make([]byte, 1024)
	n, err := hs.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(buf[:n]))
}
