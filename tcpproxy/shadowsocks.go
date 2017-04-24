package tcpproxy

import (
	"net"
	"github.com/shadowsocks/go-shadowsocks2/core"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type ShadowSocksDialer struct {
	cipher  core.Cipher //interface include StreamConn() and PacketConn()
	address string
}
