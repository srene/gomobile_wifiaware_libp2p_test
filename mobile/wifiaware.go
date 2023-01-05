package datahop

type WifiAwareClientDriver interface {
	Connect(peerId string)
	Disconnect()
	Host() string
}

type WifiAwareServerDriver interface {
	Start(peerId string, port int)
	Stop()
}

type WifiAwareNotifier interface {
	OnConnectionFailure(message string)
	OnConnectionServerSuccess(ip string, port int, peerId string)
	OnConnectionClientSuccess(listenip string, listenport int, ip string, port int, peerId string)
	OnDisconnect()
}
