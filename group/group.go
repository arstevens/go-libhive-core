package group

import (
	"github.com/arstevens/go-libhive-core/message"
	"github.com/arstevens/go-libhive-core/stream"
)

type Group struct {
	nodes   map[string]*stream.Stream
	subnets map[string][]string
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
