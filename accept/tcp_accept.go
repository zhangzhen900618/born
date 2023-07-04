package accept

import (
	"github.com/zhangzhen900618/born/logger"
)

var _ Accept = (*tcpServer)(nil)

type tcpServer struct {
	options
	connector
}

func NewTcp(address string, opts ...Option) Accept {
	o := getDefaultOptions()
	o.address = address
	for _, opt := range opts {
		opt(&o)
	}

	a := &tcpServer{
		options:   o,
		connector: newConnector(o.chanSize),
	}
	return a
}

func (s *tcpServer) Start() {
	listener, err := s.GetListener(s.certFile, s.keyFile, s.address)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	s.Run()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Err(err).Send()
				return
			}
			s.InChan(conn)
		}
	}()
}

func (s *tcpServer) Stop() {
	err := s.listener.Close()
	if err != nil {
		logger.Err(err).Send()
	}
	s.Close()
}
