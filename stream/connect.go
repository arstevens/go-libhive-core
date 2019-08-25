package stream

import (
	"context"
	"errors"
	"net"

	"github.com/arstevens/go-libhive-core/protocol"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	maddr "github.com/multiformats/go-multiaddr"
)

// ipfs api libp2p forward and listen functions
func streamListen(sh *ipfsapi.Shell, proto protocol.ID, addr maddr.Multiaddr) (*net.Listener, error) {
	resp, err := sh.Request("p2p/listen", proto.String(), addr.String()).Send(context.Background())
	if err != nil {
		return nil, err
	}

	defer resp.Close()
	if resp.Error != nil {
		return nil, resp.Error
	}

	lnProto, lnAddr, err := parseMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	ln, err := net.Listen(lnProto, lnAddr)
	return &ln, err
}

func streamForward(sh *ipfsapi.Shell, proto protocol.ID, addr maddr.Multiaddr, nid string) (*net.Conn, error) {
	resp, err := sh.Request("p2p/forward", proto.String(), addr.String(), "/ipfs/"+nid).Send(context.Background())
	if err != nil {
		return nil, err
	}

	defer resp.Close()
	if resp.Error != nil {
		return nil, resp.Error
	}

	fwProto, fwAddr, err := parseMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial(fwProto, fwAddr)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

// Attempt to discover and connect to node with id 'nid'
func establishConnection(sh *ipfsapi.Shell, nid string) error {
	// DHT find peer multiaddresses and store them in peerstore
	pi, err := sh.FindPeer(nid)
	if err != nil {
		return err
	} else if len(pi.Addrs) < 1 {
		return err
	}

	// Attempt to establish a connection to node
	return sh.SwarmConnect(context.Background(), nid)
}

// Deconstructs a multiaddr 'addr' into address 'IP:Port' and transfer protocol 'tpc/udp/etc'
func parseMultiaddr(addr maddr.Multiaddr) (string, string, error) {
	maProtos := addr.Protocols()
	if len(maProtos) < 2 {
		return "", "", errors.New("Not enough arguments in Multiaddress")
	}

	proto := maProtos[1].Name
	port, err := addr.ValueForProtocol(maProtos[1].Code)
	if err != nil {
		return "", "", err
	}

	ip, err := addr.ValueForProtocol(maProtos[0].Code)
	if err != nil {
		return "", "", err
	}
	fullAddr := ip + ":" + port
	return proto, fullAddr, nil
}

// Removes the libp2p stream listener for 'proto'
func closeProtoListener(sh *ipfsapi.Shell, proto protocol.ID) error {
	_, err := sh.Request("p2p/close", "--protocol="+proto.String()).Send(context.Background())
	return err
}
