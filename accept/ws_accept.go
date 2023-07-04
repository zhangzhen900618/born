package accept

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zhangzhen900618/born/logger"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	_ net.Conn     = (*wsConn)(nil)
	_ http.Handler = (*wsServer)(nil)
	_ Accept       = (*wsServer)(nil)
)

type (
	wsServer struct {
		options
		connector
		srv     *http.Server
		upgrade *websocket.Upgrader
	}

	wsConn struct {
		*websocket.Conn
		typ    int
		reader io.Reader
	}
)

func NewWs(address string, opts ...Option) Accept {
	o := getDefaultOptions()
	o.address = address
	for _, opt := range opts {
		opt(&o)
	}

	a := &wsServer{
		options: o,
		upgrade: &websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
		connector: newConnector(o.chanSize),
	}
	return a
}

func (s *wsServer) Start() {
	listener, err := s.GetListener(s.certFile, s.keyFile, s.address)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}
	s.Run()

	go func() {
		s.srv = &http.Server{
			Handler: s,
		}
		err = s.srv.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			logger.Err(err).Send()
		}
	}()
}

func (s *wsServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := s.srv.Shutdown(ctx)
	if err != nil {
		logger.Err(err).Send()
	}
	s.Close()
}

func (s *wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &wsConn{
		Conn: conn,
	}
	s.InChan(c)
}

func (c *wsConn) Read(b []byte) (n int, err error) {
	if c.reader == nil {
		c.typ, c.reader, err = c.NextReader()
		if err != nil {
			return n, err
		}
	}
	n, err = c.reader.Read(b)
	if err != nil && err != io.EOF {
		return n, err
	} else if err == io.EOF {
		_, c.reader, err = c.NextReader()
		if err != nil {
			return 0, err
		}
	}
	return n, err
}

func (c *wsConn) Write(b []byte) (n int, err error) {
	err = c.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return n, err
	}
	return len(b), err
}

func (c *wsConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}
