package impl

import (
	"github.com/kade-chen/google-billing-console/apps/configs"
	"github.com/kade-chen/library/ioc"
)

// var _ configs.Service = (*Service)(nil)

func init() {
	ioc.Config().Registry(&Service{})
}

type Service struct {
	ioc.ObjectImpl
	// log                          *zerolog.Logger
	Default_Project_ID           string `toml:"default_project_id" json:"default_project_id" yaml:"default_project_id"`
	Default_Service_Account_Name string `toml:"default_service_account_name" json:"default_service_account_name" yaml:"default_service_account_name"`
}

func (s *Service) Init() error {
	return nil
}

func (Service) Name() string {
	return configs.AppName
}

func (i *Service) Priority() int {
	return 0
}
