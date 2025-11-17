package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/domain"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps/domain/impl"
)

var (
	ctx  context.Context
	impl domain.Service
)

func TestCreateDomain(t *testing.T) {
	req := domain.NewCreateDomainRequest()
	req.Name = "wondercloud.com"
	req.Description = "test domain"
	ins, _ := impl.CreateDomain(ctx, req)
	fmt.Println(format.ToJSON(ins))
}

func TestDescribeDomain(t *testing.T) {
	// req := domain.NewDescribeDomainRequestByName(domain.DEFAULT_DOMAIN)
	req := domain.NewDescribeDomainRequestByName("wondercloud.com")
	ins, err := impl.DescribeDomain(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(ins)
	fmt.Println(format.ToJSON(ins))
	// t.Log(ins.ToJson())
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go-kade-project/github/google-billing-console/etc/config.toml"
	ioc.DevelopmentSetup(req)
	impl = ioc.Controller().Get(domain.AppName).(domain.Service)
}
