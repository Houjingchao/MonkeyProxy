package tcpproxy

import (
	"net"
)

type Client struct {
	Address string
	Target  string
	Dialer  Dialer
	i       int64
}

func NewClient(address, target string, dialer Dialer) *Client { // new a client
	return &Client{Address: address, Target: target, Dialer: dialer}
}

func (client *Client) run() {
	server, err := net.Listen("tcp", client.Address)
	if err != nil {
		panic(err) //execution stop immediately
	}

	for {
		if conn, err := server.Accept(); err == nil { // execute right
			//do something
			go client.handle(conn)
		}
	}
}

func (client *Client) handle(conn net.Conn) {

}
