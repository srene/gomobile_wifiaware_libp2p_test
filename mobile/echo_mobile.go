package datahop

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"

	//	logger "github.com/ipfs/go-log/v2"

	"github.com/srene/go-libp2p"
	"github.com/srene/go-libp2p/core/crypto"
	"github.com/srene/go-libp2p/core/host"
	"github.com/srene/go-libp2p/core/network"
	"github.com/srene/go-libp2p/core/peer"
	"github.com/srene/go-libp2p/core/peerstore"

	ma "github.com/srene/go-multiaddr"
)

var (
	test *echo
)

type echo struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wService   *wifiAwareService
	ha         host.Host
	listenPort int
}

func Init(listenPort int) error {

	ctx, cancel := context.WithCancel(context.Background())

	test = &echo{
		ctx:        ctx,
		cancel:     cancel,
		listenPort: listenPort,
	}

	service, err := NewWifiAwareService()
	if err != nil {
		fmt.Println("ble discovery setup failed : ", err.Error())
		cancel()
		return err
	}
	if res, ok := service.(*wifiAwareService); ok {
		test.wService = res
	}
	return nil

}

func StartListener() string {

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	ha, err := makeBasicHost("::", test.listenPort, false, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ha.ID())

	test.ha = ha
	fmt.Printf("Starting Listener")

	fullAddr := getHostAddress(test.ha)
	fmt.Printf("I am %s\n", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.
	test.ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		fmt.Println("listener received new stream")
		if err := doEcho(s); err != nil {
			fmt.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	fmt.Println("listening for connections")

	return test.ha.ID().String()
}

func RunSender(targetPeer string) {

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	fmt.Printf("Starting Sender")

	fullAddr := getHostAddress(test.ha)
	fmt.Printf("I am %s\n", fullAddr)

	// Set a stream handler on host A. /echo/1.0.0 is
	// a user-defined protocol name.

	test.ha.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		fmt.Println("sender received new stream")
		if err := doEcho(s); err != nil {
			fmt.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})

	fmt.Printf("Connecting to %s\n", targetPeer)

	// Turn the targetPeer into a multiaddr.
	maddr, err := ma.NewMultiaddr(targetPeer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// We have a peer ID and a targetAddr so we add it to the peerstore
	// so LibP2P knows how to contact it
	test.ha.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	fmt.Printf("sender opening stream %s %s", info.ID, info.Addrs)
	// make a new stream from host B to host A
	// it should be handled on host A by the handler we set above because
	// we use the same /echo/1.0.0 protocol
	s, err := test.ha.NewStream(context.Background(), info.ID, "/echo/1.0.0")
	if err != nil {

		fmt.Printf("New stream error %s", err)
		return
	}

	fmt.Println("sender saying hello")
	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := io.ReadAll(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("read reply: %q\n", out)
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It won't encrypt the connection if insecure is true.
func makeBasicHost(listenAddress string, listenPort int, insecure bool, randseed int64) (host.Host, error) {
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip6/%s/tcp/%d", listenAddress, listenPort)),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}

	if insecure {
		opts = append(opts, libp2p.NoSecurity)
	}

	return libp2p.New(opts...)
}

func getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}

// GetDiscoveryNotifier returns discovery notifier
func GetPeerId() string {
	return test.ha.ID().String()
}

// GetDiscoveryNotifier returns discovery notifier
func GetWifiAwareNotifier() WifiAwareNotifier {
	return test.wService
}
