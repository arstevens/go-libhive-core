package group

import (
	"strings"

	"github.com/arstevens/go-libhive-core/message"
	"github.com/arstevens/go-libhive-core/stream"
)

type Group struct {
	nodes     map[string]*stream.Stream
	subnets   map[string][]string
	entryNode string
	exitNode  string
}

func NewGroup(n map[string]*stream.Stream) *Group {
	g := Group{
		nodes:     n,
		subnets:   make(map[string][]string),
		entryNode: "",
		exitNode:  "",
	}

	g.evalSubnet()
	return &g
}

func (g *Group) Add(nid string, s *stream.Stream) {
	g.nodes[nid] = s
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
	return g.entryNode, g.nodes[g.entryNode]
}

func (g *Group) ExitNode() (string, *stream.Stream) {
	return g.exitNode, g.nodes[g.exitNode]
}

func (g *Group) evalSubnet() {
	for k, _ := range g.nodes {
		if g.entryNode == "" {
			g.entryNode = k
			g.exitNode = k
		} else if strings.Compare(k, g.entryNode) > 0 {
			g.entryNode = k
		}

		if strings.Compare(k, g.exitNode) < 0 {
			g.exitNode = k
		}
	}
}
