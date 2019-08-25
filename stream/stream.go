package stream

import (
	"context"
	"net"
	"strconv"

	"github.com/arstevens/go-libhive-core/protocol"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/phayes/freeport"
)

// Stream type and reciever methods
// Stream wraps a net.Conn for readability
type Stream struct {
	conn     *net.Conn
	sHandler bool
}

func wrapConn(c *net.Conn, handler bool) Stream {
	return Stream{conn: c, sHandler: handler}
}

// Private Accessors
func (s *Stream) getConn() *net.Conn { return s.conn }
func (s *Stream) isHandler() bool    { return s.sHandler }

// Implementing io.ReadWriter interface
func (s *Stream) Read(b []byte) (int, error) {
	c := *s.getConn()
	return c.Read(b)
}

func (s *Stream) Write(b []byte) (int, error) {
	c := *s.getConn()
	return c.Write(b)
}

func (s *Stream) Close() error {
	if !s.isHandler() {
		// close ipfs p2p forwarder
	}

	c := *s.getConn()
	err := c.Close()

}

func NewStream(sh *ipfsapi.Shell, proto protocol.ID, id string) (*Stream, error) {

}

// NewStreamHandler should be run in its own go routine
// NSH runs 'callback' on every new stream that connects to 'proto'
func NewStreamHandler(sh *ipfsapi.Shell, proto protocol.ID, callback func(s Stream)) {
	fport, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	addr, err := maddr.NewMultiaddr("/ip4/127.0.0.1/tcp/" + strconv.Itoa(fport))
	if err != nil {
		panic(err)
	}

	ln, err := streamListen(sh, proto, addr)
	if err != nil {
		panic(err)
	}
	// defer closing p2p listener

	for {
		conn, err := (*ln).Accept()
		if err != nil {
			continue
		}
		s := wrapConn(&conn, true)
		callback(s)
		s.Close()
	}
}

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

	maProtos := addr.Protocols()
	if len(maProtos) < 2 {
		return nil, err
	}

	lnProto := maProtos[1].Name
	lnPort, err := addr.ValueForProtocol(maProtos[1].Code)
	if err != nil {
		return nil, err
	}

	ip, err := addr.ValueForProtocol(maProtos[0].Code)
	if err != nil {
		return nil, err
	}
	lnAddr := ip + ":" + lnPort

	ln, err := net.Listen(lnProto, lnAddr)
	return &ln, err
}

func streamForward(sh *ipfsapi.Shell, proto protocol.ID, addr maddr.Multiaddr, nid string) (*net.Conn, error) {

}
