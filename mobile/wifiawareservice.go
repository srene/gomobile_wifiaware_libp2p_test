package datahop

import (
	"fmt"
)

type wifiAwareService struct {
	peerId string
}

func NewWifiAwareService(peerId string) (WifiAwareNotifier, error) {

	waService := &wifiAwareService{
		peerId: peerId,
	}
	return waService, nil
}

func (w *wifiAwareService) OnConnectionFailure(message string) {

}
func (w *wifiAwareService) OnConnectionServerSuccess(ip string, port int) {
	fmt.Printf("Connection server success  %s port %d\n", ip, port)
	StartListener()
	//RunSender("::", port+1, fmt.Sprintf("/ip6/%s/tcp/%d/p2p/%s", ip, port, peerId))

}
func (w *wifiAwareService) OnConnectionClientSuccess(ip string, zone string, port int, peerId string) {
	fmt.Printf("Connection client success peerid %s %s port %d\n", peerId, ip, port)
	RunSender(fmt.Sprintf("/ip6zone/%s/ip6/%s/tcp/%d/p2p/%s", zone, ip, port, peerId))
}
func (w *wifiAwareService) OnDisconnect() {

}
