package datahop

import (
	"fmt"
	"strconv"
	"time"
)

type wifiAwareService struct {
	peerId string
}

func NewWifiAwareService() (WifiAwareNotifier, error) {

	waService := &wifiAwareService{}
	return waService, nil
}

func (w *wifiAwareService) OnConnectionFailure(message string) {

}
func (w *wifiAwareService) OnConnectionServerSuccess(ip string, zone string, port int) string {
	//address := "[" + ip + "%" + zone + "]:" + strconv.Itoa(port)
	address := "[::]:" + strconv.Itoa(port)

	fmt.Printf("Connection server success  %s\n", address)
	//time.Sleep(time.Second * 5)
	//return StartListener(fmt.Sprintf("/ip6/%s/tcp/%d", "::", port))
	go StartListener(address)
	//RunSender("::", port+1, fmt.Sprintf("/ip6/%s/tcp/%d/p2p/%s", ip, port, peerId))
	return "QmVoW6bbThjrahCu8AdkLewnXhLp2nDHxNFDVudMZiCFQn"
}
func (w *wifiAwareService) OnConnectionClientSuccess(ip string, zone string, port int, peerId string) {
	//fmt.Printf("Connection client success peerid %s %s port %d\n", peerId, ip, port)
	//StartListener(fmt.Sprintf("/ip6/::/tcp/%d", port+1))
	address := "[" + ip + "%" + zone + "]:" + strconv.Itoa(port)
	//address := "[::1]:" + strconv.Itoa(port)

	//laddress := "[::]:" + strconv.Itoa(port+1)
	fmt.Printf("Connection client success peerid %s\n", address)
	time.Sleep(time.Second * 2)
	//StartListener(laddress)
	//RunSender(fmt.Sprintf("/ip6zone/%s/ip6/%s/tcp/%d/p2p/%s", zone, ip, port, peerId))
	go RunSender(address)

}
func (w *wifiAwareService) OnDisconnect() {

}
