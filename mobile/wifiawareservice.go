package datahop

import (
	"fmt"
)

type wifiAwareService struct {
	peerId string
}

func NewWifiAwareService() (WifiAwareNotifier, error) {

	waService := &wifiAwareService{}
	return waService, nil
}

func (w *wifiAwareService) Start() {

}

func (w *wifiAwareService) Stop() {

}

func (w *wifiAwareService) OnConnectionFailure(message string) {

}
func (w *wifiAwareService) OnConnectionServerSuccess(ip string, port int, peerId string) {
	fmt.Printf("Connection server success peerid %s %s port %d\n", peerId, ip, port)
	//StartListener(ip, port)
	//RunSender("::", port+1, fmt.Sprintf("/ip6/%s/tcp/%d/p2p/%s", ip, port, peerId))

}
func (w *wifiAwareService) OnConnectionClientSuccess(listenip string, listenport int, ip string, port int, peerId string) {
	fmt.Printf("Connection client success peerid %s %s port %d\n", peerId, ip, port)
	//RunSender(listenip, listenport, fmt.Sprintf("/ip6zone/%s/tcp/%d/p2p/%s", ip, port, peerId))
}
func (w *wifiAwareService) OnDisconnect() {

}
