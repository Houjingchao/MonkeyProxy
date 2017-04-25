package MonkeyProxy

import (
	"MonkeyProxy/configmanager"
	"github.com/lextoumbourou/goodhosts"
	"fmt"
	"log"
	"net"
	"github.com/Unknwon/com"
	"MonkeyProxy/tcpproxy"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//get config
	config := configmanager.MustLoad()
	hosts, _ := goodhosts.NewHosts()

	fmt.Println("write /etc/hosts")
	defer hosts.Flush()

	for _, v, err := range config.ProxyTargets {
		if err != nil {
			log.Fatal(err)
		}
		listen, _, err := net.SplitHostPort(v.Listen)
		if err != nil {
			log.Fatal(err)
		}
		host := listen
		if v.Host != "" {
			host = v.Host
		}
		log.Printf("# 添加 HOST %v %v", listen, host)

		if !hosts.Has(listen, host) {
			hosts.Add(listen, host)
		}
		log.Printf("# 添加 lo0 alias %v", listen)

		stdout, stderr, err := com.ExecCmd("ifconfig", "lo0", "alias", listen, "up")
		if err != nil {
			fmt.Errorf("Fail to clone docs from remote source(%s): %v - %s", listen, err, stderr)
		}
		fmt.Println(stdout)
	}

	if err := hosts.Flush(); err != nil {
		log.Println("please sudo run")
	}

	server := config.Server
	dialer := tcpproxy.NewShadowSocksDialer(server.Address, server.CipherType, server.Password)

	for _, v := range config.ProxyTargets {
		client := tcpproxy.NewClient(v.Listen, v.Target, dialer)
		host := v.Listen
		if v.Host != "" {
			host = v.Host
			log.Printf("代理 %v <-> %v <-> %v", host, client.Address, client.Target)
		} else {
			log.Printf("代理 %v <-> %v <-> %v", client.Target, client.Address, client.Target)
		}
		go client.Run()
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh // block here
	fmt.Println("# kill me  success")
}
