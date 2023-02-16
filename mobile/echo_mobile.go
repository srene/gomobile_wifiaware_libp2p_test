package datahop

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"

	//	logger "github.com/ipfs/go-log/v2"

	"github.com/libp2p/go-libp2p/core/host"

	"github.com/ipfs/go-log/v2"
)

var (
	test *echo
)

const (
	CONN_TYPE = "tcp6"
)

type echo struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wService   *wifiAwareService
	ha         host.Host
	listenPort int
}

func Init(listenPort int) error {

	log.SetAllLoggers(log.LevelDebug)

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

func StartListener(address string) {

	fmt.Println("Starting listener " + address)

	l, err := net.Listen(CONN_TYPE, address)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		//os.Exit(1)
		return
	}
	defer l.Close()
	fmt.Println("Listening on " + address)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			//os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func RunSender(address string) {

	fmt.Println("Connecting " + CONN_TYPE + " " + address)
	conn, err := net.Dial(CONN_TYPE, address)
	if err != nil {
		fmt.Println("Error connecting " + err.Error())
		return
		//os.Exit(1)
	}

	ReadNWrite(conn)
	conn.Close()
}

// GetDiscoveryNotifier returns discovery notifier
func GetPeerId() string {
	return "QmVoW6bbThjrahCu8AdkLewnXhLp2nDHxNFDVudMZiCFQn"
}

// GetDiscoveryNotifier returns discovery notifier
func GetWifiAwareNotifier() WifiAwareNotifier {
	return test.wService
}

func ReadNWrite(conn net.Conn) {
	fmt.Println("TEst request")
	message := "Test Request\n"
	_, write_err := conn.Write([]byte(message))
	if write_err != nil {
		fmt.Println("failed:", write_err)
		return
	}
	conn.(*net.TCPConn).CloseWrite()

	buf, read_err := ioutil.ReadAll(conn)
	if read_err != nil {
		fmt.Println("failed:", read_err)
		return
	}
	fmt.Println(string(buf))
}

func handleRequest(conn net.Conn) {
	buf, read_err := ioutil.ReadAll(conn)
	if read_err != nil {
		fmt.Println("failed:", read_err)
		return
	}
	fmt.Println("Got: ", string(buf))

	_, write_err := conn.Write([]byte("Message received.\n"))
	if write_err != nil {
		fmt.Println("failed:", write_err)
		return
	}
	conn.Close()
}
