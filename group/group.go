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

func (g *Group) Publish(data message.Message) error {
	buf := make([]byte, 8192)
	for k, v := range g.nodes {
		n, err := data.Read(buf)
		if n > 0 {
			_, err := v.Write(buf[:n])
		}
		v.Write(data)
		_, err := v.Write
	}
}
