package logger_test

import (
	"github.com/zhangzhen900618/born/logger"
	"testing"
)

func TestDebug(t *testing.T) {
	logger.New(logger.WithEnableWrite(true))
	logger.Debug().Str("name", "test").Send()
}
