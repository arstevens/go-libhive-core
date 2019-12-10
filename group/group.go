package group

import (
	"strings"

	"github.com/arstevens/go-libhive-core/message"
	"github.com/arstevens/go-libhive-core/stream"
	ipfsapi "github.com/ipfs/go-ipfs-api"
)

type Group struct {
	nodes     map[string]*stream.Stream
	sortedNodeKeys []string
	subnets   map[string][]string
	ipfsShell *ipfsapi.Shell
}

func NewGroup(sh *ipfsapi.Shell, n map[string]*stream.Stream) *Group {
	g := Group{
		nodes:     n,
		subnets:   make(map[string][]string),
		ipfsShell: sh,
	}

	nids := make([]string, len(g.nodes))
	i := 0
	for k, _ := range g.nodes {
		nids[i] = k
		i += 1
	}
	sort.Strings(nids)
	g.sortedNodeKeys = nids

	return &g
}

func (g *Group) Add(nid string, s *stream.Stream) {
	g.nodes[nid] = s
	g.sortedNodeKeys = append(g.sortedNodeKeys, nid)
	sort.Strings(g.sortedNodeKeys)
}

func (g *Group) SortedKeys() []string {
	return g.sortedNodeKeys
}

func (g *Group) Tag(nid string, subid string) {
	subnet, ok := g.subnets[subid]
	if !ok {
		subnet = make([]string, 1)
		subnet[1] = nid
	} else {
		subnet = append(subnet, nid)
	}
	g.subnets[subid] = subnet
}

func (g *Group) GetNodes() map[string]*stream.Stream {
	return g.nodes
}

func (g *Group) GetShell() *ipfsapi.Shell {
	return g.ipfsShell
}

func (g *Group) publish(data message.Message, nodes map[string]*stream.Stream) error {
	buf := make([]byte, 8192)
	for _, v := range nodes {
		n, _ := data.Read(buf)
		for n > 0 {
			_, err := v.Write(buf[:n])
			if err != nil {
				return err
			}
			n, err = data.Read(buf)
		}
		data.Reset()
	}
	return nil
}

func (g *Group) Publish(data message.Message) error {
	return g.publish(data, g.nodes)
}

func (g *Group) TagPublish(data message.Message, tag string) error {
	subnetNodeIds := g.subnets[tag]
	subnetMap := make(map[string]*stream.Stream)
	for _, node := range subnetNodeIds {
		subnetMap[node] = g.nodes[node]
	}
	return g.publish(data, subnetMap)
}

func (g *Group) EntryNode() (string, *stream.Stream) {
	nid := g.sortedNodeKeys[0]
	return nid, g.nodes[nid]
}

func (g *Group) ExitNode() (string, *stream.Stream) {
	nid := g.sortedNodeKeys[len(g.sortedNodeKeys) - 1]
	return nid, g.nodes[nid]
}
