package cache_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/cache"
	_ "github.com/kade-chen/library/ioc/config/cache"
	_ "github.com/kade-chen/library/ioc/config/log"
)

func TestObjectLoad(t *testing.T) {
	fmt.Println("shhshsh")
	a := ioc.Config().Get(cache.AppName).(*cache.Cache)
	m := a.Ristretto.Metrics
	log.Printf("hit=%d miss=%d reject=%d",
		m.Hits(), m.Misses(), m.SetsRejected(),
	)
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/library/etc/application.toml"
	ioc.DevelopmentSetup(req)
}
