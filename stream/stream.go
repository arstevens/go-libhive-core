package stream

import (
	"context"
	"fmt"
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
	sh       *ipfsapi.Shell
	proto    protocol.ID
}

func wrapConn(c *net.Conn, handler bool, sh *ipfsapi.Shell, proto protocol.ID) Stream {
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
		(*s).sh.Request("p2p/close", "--protocol="+(*s).proto.String()).Send(context.Background())
	}

	c := *s.getConn()
	return c.Close()
}

func NewStream(sh *ipfsapi.Shell, proto protocol.ID, nid string) (*Stream, error) {
	fmt.Println("Attempting to establish connection")
	err := establishConnection(sh, nid)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connection Established")

	fport, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}
	fmt.Println("Port Allocated")

	addr, err := maddr.NewMultiaddr("/ip4/127.0.0.1/tcp/" + strconv.Itoa(fport))
	if err != nil {
		return nil, err
	}

	conn, err := streamForward(sh, proto, addr, nid)
	if err != nil {
		return nil, err
	}
	fmt.Println("Stream forwarding")

	s := wrapConn(conn, false, sh, proto)
	fmt.Println("Wrapping connection object")
	return &s, nil
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
	defer closeProtoListener(sh, proto)

	for {
		conn, err := (*ln).Accept()
		if err != nil {
			continue
		}
		s := wrapConn(&conn, true, sh, proto)
		callback(s)
		s.Close()
	}
}
