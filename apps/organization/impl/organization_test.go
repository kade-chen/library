package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/google-billing-console/apps/organization"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/tools/format"

	_ "github.com/kade-chen/google-billing-console/apps/configs/impl"
	_ "github.com/kade-chen/google-billing-console/apps/organization/impl"
)

var (
	ctx  = context.Background()
	impl organization.Service
)

func TestCreateOrganization(t *testing.T) {
	req := organization.NewCreateOrganizationRequest()
	req.OrganizationDetail.SubOrganization = "test-wondercloud.com"
	req.Description = "test Organization"
	ins, err := impl.CreateOrganization(ctx, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(format.ToJSON(ins))
}

func TestDescribeOrganization(t *testing.T) {
	// req := Organization.NewDescribeOrganizationRequestByName(Organization.DEFAULT_Organization)
	req := organization.NewDescribeOrganizationRequestByName("test-wondercloud.com")
	ins, err := impl.DescribeOrganization(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(ins)
	fmt.Println(format.ToJSON(ins))
	// t.Log(ins.ToJson())
}

func TestListOrganizations(t *testing.T) {
	// req := Organization.NewDescribeOrganizationRequestByName(Organization.DEFAULT_Organization)
	req := organization.NewListOrganizationRequest(&organization.ListOrganizationRequest{
		Page:  &organization.PageRequest{},
		Names: []string{"wondercloud.com", "test-wondercloud.com"},
	})
	//总数
	// req.Page.PageSize = 10
	// //跳过多少个数据wa
	// req.Page.Offset = 0
	ins, err := impl.ListOrganizations(ctx, req)
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
	impl = ioc.Controller().Get(organization.AppName).(organization.Service)
}
