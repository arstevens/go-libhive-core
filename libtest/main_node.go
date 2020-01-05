package main

import (
	"fmt"
	"sort"

	"github.com/arstevens/go-libhive-core/node"
)

func printNodes(n node.Nodes) {
	for _, v := range n {
		fmt.Println(v.Id())
	}
	fmt.Println("-------------------")
}

func main() {
	n1 := node.NewRemoteNode("QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa", nil, nil)
	n2 := node.NewRemoteNode("QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM", nil, nil)
	n3 := node.NewRemoteNode("QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd", nil, nil)
	n4 := node.NewRemoteNode("QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt", nil, nil)
	ln := node.NewLocalNode("QmadRNcr9cdaqaDFCmj7VfSTzUN6pwZFCs7952bRTnE5qF", nil)

	var nodes node.Nodes
	nodes = append(nodes, n3, n2, n4, n1)
	printNodes(nodes)
	sort.Sort(nodes)
	printNodes(nodes)
	l, r := nodes.FindNeighbors(ln)
	fmt.Println((*l).Id() + " " + (*r).Id())

	nodes = append(nodes, ln)
	sort.Sort(nodes)
	printNodes(nodes)
}
