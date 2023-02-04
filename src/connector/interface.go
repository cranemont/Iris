package connector

type Connector interface {
	Connect()
	Disconnect()
}
