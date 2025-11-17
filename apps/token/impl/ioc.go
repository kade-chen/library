package impl

import (
	"context"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/google-billing-console/apps/token/provider"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/grpc"
	"github.com/kade-chen/library/ioc/config/log"
	ioc_mongo "github.com/kade-chen/library/ioc/config/mongo"
)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	col *mongo.Collection
	token.UnimplementedRPCServer
	ioc.ObjectImpl
	log *zerolog.Logger

	// policy  policy.Service
	// ns      namespace.Service
	// checker security.Checker
	// domain  domain.Service
	// notify  notify.Service
}

func (s *service) Init() error {
	dc := ioc_mongo.DB().Collection(s.Name())
	indexs := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "refresh_token", Value: -1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "issue_at", Value: -1}},
		},
	}

	_, err := dc.Indexes().CreateMany(context.Background(), indexs)
	if err != nil {
		return err
	}

	s.col = dc

	s.log = log.Sub(s.Name())
	// s.ns = ioc.Controller().Get(namespace.AppName).(namespace.Service)
	// s.policy = ioc.Controller().Get(policy.AppName).(policy.Service)
	// s.domain = ioc.Controller().Get(domain.AppName).(domain.Service)
	// s.notify = ioc.Controller().Get(notify.AppName).(notify.Service)

	// s.checker, err = security.NewChecker()
	// if err != nil {
	// 	return fmt.Errorf("new checker error, %s", err)
	// }

	// 初始化所有的auth provider
	if err := provider.Init(); err != nil {
		return err
	}

	token.RegisterRPCServer(grpc.Get().Server(), s)
	return nil
}

func (service) Name() string {
	return token.AppName
}
