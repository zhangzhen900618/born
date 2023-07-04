package accept_test

import (
	"fmt"
	"github.com/zhangzhen900618/born/accept"
	"github.com/zhangzhen900618/born/logger"
	"net"
	"testing"
	"time"
)

func TestNewTcp(t *testing.T) {
	s := accept.NewTcp(":9071")
	s.OnConnect(func(conn net.Conn) {
		logger.Info().Str("new conn", conn.RemoteAddr().String()).Send()
		go func() {
			for {
				buf := make([]byte, 2048)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						return
					}
					fmt.Println(buf[:n])
				}
			}
		}()
	})
	s.Start()

	time.Sleep(10 * time.Second)
	s.Stop()

	time.Sleep(time.Second)
}
