package tcpproxy

import (
	"net"
	"github.com/shadowsocks/go-shadowsocks2/core"
	"log"
	"github.com/shadowsocks/go-shadowsocks2/socks"
	"qiniupkg.com/x/errors.v7"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type ShadowSocksDialer struct {
	cipher  core.Cipher //interface include StreamConn() and PacketConn()  cipher type
	address string
}

func NewShadowSocksDialer(address, password, cipher string) (*ShadowSocksDialer) {

	c, err := core.PickCipher(cipher, []byte{}, password) //give a cipher name from password
	if err != nil {
		log.Fatal(err)
	}
	return &ShadowSocksDialer{cipher: c, address: address}
}

func (s *ShadowSocksDialer) Dial(network, address string) (net.Conn, error) {

	conn, err := core.Dial(network, s.address, s.cipher) //core code

	if err != nil {
		return nil, err
	}
	tgt := socks.ParseAddr(address)

	if tgt == nil {
		return nil, errors.New("wrong des address")
	}
	conn.Write(tgt)
	return conn, nil
}
