package tcpproxy

import (
	"net"
	"log"
	"sync/atomic"
	"io"
	"time"
)

type Client struct {
	Address string //listen address
	Target  string // target address
	Dialer  Dialer
	i       int64
}

func NewClient(address, target string, dialer Dialer) *Client { // new a client
	return &Client{Address: address, Target: target, Dialer: dialer}
}

func (client *Client) Run() {
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

//return dst to src num,src to dst num ,err
func (client *Client) rely(dst, src net.Conn) (int64, int64, error) {

	type result struct {
		N   int64
		Err error
	}
	ch := make(chan result)

	go func() {
		n, err := io.Copy(dst, src) //copy src to dst ,return number of the bytes
		src.SetDeadline(time.Now())
		dst.SetDeadline(time.Now())
		ch <- result{n, err}
	}()

	copynumber, err := io.Copy(dst, src) //set deadline associated with connection
	src.SetDeadline(time.Now())
	dst.SetDeadline(time.Now())

	rs := <-ch

	if err == nil { // if not, update the err
		err = rs.Err
	}
	return copynumber, rs.N, err
}
func (client *Client) handle(conn net.Conn) {

	defer conn.Close()
	// core code
	targetConn, err := client.Dialer.Dial("tcp", client.Target)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer targetConn.Close()

	i := atomic.AddInt64(&client.i, 1)
	log.Println("begin to switch:", client.Target, i)
	_, _, err = client.rely(targetConn, conn)

	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Println("switch end~")
		}
		log.Fatalf("switch failed:%v %d", err, i)
	}
}
