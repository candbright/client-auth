package client

import "github.com/go-resty/resty/v2"

var globalClient Client

func Ins() Client {
	if globalClient == nil {
		globalClient = Default()
	}
	return globalClient
}

func Init(endpoint string) {
	globalClient = New(endpoint)
}

func Default() Client {
	return New("http://localhost:11100")
}

func New(endpoint string) Client {
	return &client{endpoint: endpoint, client: resty.New()}
}
