package accept

import (
	"crypto/tls"
	"net"
)

type (
	OnConnect func(conn net.Conn)

	Accept interface {
		Start()
		Stop()
		OnConnect(fn OnConnect)
	}

	connector struct {
		listener  net.Listener
		onConnect OnConnect
		connChan  chan net.Conn
	}
)

func newConnector(size int) connector {
	c := connector{
		connChan: make(chan net.Conn, size),
	}
	return c
}

func (c *connector) OnConnect(fn OnConnect) {
	c.onConnect = fn
}

func (c *connector) InChan(conn net.Conn) {
	c.connChan <- conn
}

func (c *connector) Run() {
	go func() {
		for conn := range c.connChan {
			c.onConnect(conn)
		}
	}()
}

func (c *connector) Close() {
	close(c.connChan)
}

func (c *connector) GetListener(certFile, keyFile, address string) (net.Listener, error) {
	var err error
	if certFile == "" || keyFile == "" {
		c.listener, err = net.Listen("tcp", address)
		return c.listener, err
	}

	crt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{crt},
	}

	c.listener, err = tls.Listen("tcp", address, tlsCfg)
	return c.listener, err
}
