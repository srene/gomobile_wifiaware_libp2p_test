package datahop

type WifiAwareClientDriver interface {
	Connect(peerId string)
	Disconnect()
}

type WifiAwareServerDriver interface {
	Start(peerId string, port int)
	Stop()
}

type WifiAwareNotifier interface {
	OnConnectionFailure(message string)
	OnConnectionServerSuccess(ip string, port int) string
	OnConnectionClientSuccess(ip string, zone string, port int, peerId string)
	OnDisconnect()
}
