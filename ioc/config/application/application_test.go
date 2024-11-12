package application_test

import (
	"context"
	"os"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/application"
	"github.com/kade-chen/library/ioc/server"
	"github.com/BurntSushi/toml"
)

func TestDefaultConfig(t *testing.T) {
	f, err := os.OpenFile("test/default.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	appConf := map[string]any{application.AppName: application.Get()}
	toml.NewEncoder(f).Encode(appConf)
	server.Run(context.Background(), ioc.NewLoadConfigRequest())
}

func init() {
	os.Setenv("HTTP_ENABLE_TRACE", "false")
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "test/application.toml"
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}
