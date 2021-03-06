package stream

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

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
	proto    string
}

func wrapConn(c *net.Conn, handler bool, sh *ipfsapi.Shell, proto string) Stream {
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

func (s *Stream) ChanReader() chan []byte {
	ch := make(chan []byte)
	go connListener(s.getConn(), ch)
	return ch
}

func (s *Stream) Write(b []byte) (int, error) {
	c := *s.getConn()
	return c.Write(b)
}

func (s *Stream) WriteReader(r io.Reader) error {
	c := *s.getConn()
	buf := make([]byte, 8192)

	n, _ := r.Read(buf)
	for n > 0 {
		_, err := c.Write(buf[:n])
		if err != nil {
			fmt.Println(err.Error())
		}
		n, _ = r.Read(buf)
	}
	return nil
}

func (s *Stream) Close() error {
	if !s.isHandler() {
		(*s).sh.Request("p2p/close", "--protocol="+(*s).proto).Send(context.Background())
	}

	c := *s.getConn()
	return c.Close()
}

func NewStream(sh *ipfsapi.Shell, proto string, nid string) (*Stream, error) {
	// Will be replaced by WDS multiaddress cacheing functionality
	/*
		fmt.Println("Attempting to establish connection")
		err := establishConnection(sh, nid)
		if err != nil {
			return nil, err
		}
		fmt.Println("Connection Established")
	*/

	fport, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}

	addr, err := maddr.NewMultiaddr("/ip4/127.0.0.1/tcp/" + strconv.Itoa(fport))
	if err != nil {
		return nil, err
	}

	conn, err := streamForward(sh, proto, addr, nid)
	if err != nil {
		return nil, err
	}

	s := wrapConn(conn, false, sh, proto)
	return &s, nil
}

// NewStreamHandler should be run in its own go routine
// NSH runs 'callback' on every new stream that connects to 'proto'
func NewStreamHandler(sh *ipfsapi.Shell, proto string, callback func(s Stream)) {
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

func connListener(c *net.Conn, ch chan []byte) {
	// find optimal sleep time
	buf := make([]byte, 8192)
	for {
		n, _ := (*c).Read(buf)
		for n > 0 {
			ch <- buf[:n]
			var err error
			n, _ = (*c).Read(buf)
			if err != nil {
				return
			}
		}
		time.Sleep(time.Second / 2)
	}
}
