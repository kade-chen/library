package log_test

import (
	"os"
	"runtime/debug"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/BurntSushi/toml"
)

func TestGetClientGetter(t *testing.T) {
	sub := log.Sub("module_p")
	// log.T("module_a").Trace(context.Background())
	sub.Debug().Msgf("hello %s", "a")
}

func TestDefaultConfig(t *testing.T) {
	f, err := os.OpenFile("test/default.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	appConf := map[string]any{log.AppName: ioc.Config().Get(log.AppName).(*log.Config)}
	toml.NewEncoder(f).Encode(appConf)
	sub := log.Sub("module_p")
	// log.T("module_a").Trace(context.Background())
	sub.Debug().Msgf("hello %s", "a")
}

func TestPanicStack(t *testing.T) {
	// 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			log.L().Error().Stack().Msgf("Panic occurred: %v\n%s", err, debug.Stack())
		}
	}()

	// 代码中可能发生 panic 的地方
	panic("Something went wrong!")
}

func init() {
	a := ioc.Config().Get(log.AppName).(*log.Config)
	a.File.Enable = true
	a.File.FilePath = "/Users/kade.chen/go12-project/library/ioc/config/log/logs/2.log"
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
