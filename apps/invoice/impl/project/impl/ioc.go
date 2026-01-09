package impl

import (
	"github.com/rs/zerolog"

	"cloud.google.com/go/bigquery"

	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/google-billing-console/apps/configs/impl"
	"github.com/kade-chen/google-billing-console/apps/invoice"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/project"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/services"
	"github.com/kade-chen/google-billing-console/apps/invoice/impl/sku"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
)

var _ invoice.ProjectService = (*service)(nil)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	// col *mongo.Collection
	// token.UnimplementedRPCServer
	ioc.ObjectImpl
	log  *zerolog.Logger
	bq   *bigquery.Client
	svcs invoice.Service
	skus invoice.SkuService

	// policy  policy.Service
	// ns      namespace.Service
	// checker security.Checker
	// organization  organization.Service
	// notify  notify.Service
}

func (s *service) Init() error {
	s.log = log.Sub(s.Name())
	s.bq = ioc.Config().Get(configs.AppName).(*impl.Service).BQ

	s.svcs = ioc.Controller().Get(services.AppName).(invoice.Service)
	s.skus = ioc.Controller().Get(sku.AppName).(invoice.SkuService)
	s.log.Debug().Msgf("%s init successful", s.Name())
	return nil
}

func (service) Name() string {
	return project.AppName
}

func (i *service) Priority() int {
	return 0
}
