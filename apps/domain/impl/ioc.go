package impl

import (
	"context"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/grpc"
	logs "github.com/kade-chen/library/ioc/config/log"
	mongodb "github.com/kade-chen/library/ioc/config/mongo"

	"github.com/kade-chen/google-billing-console/apps/domain"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ domain.Service = (*service)(nil)

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	col *mongo.Collection
	domain.UnimplementedRPCServer
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (*service) Name() string {
	return domain.AppName
}

func (s *service) Init() error {

	s.log = logs.Sub(domain.AppName)
	s.col = mongodb.DB().Collection(domain.AppName)

	indexs := func() []mongo.IndexModel {
		var Unique bool = true
		indexs := []mongo.IndexModel{
			{Keys: bson.D{{Key: "create_at", Value: 1}}},
			{
				Keys: bson.D{
					{Key: "name", Value: -1},
				},
				Options: &options.IndexOptions{Unique: &Unique},
			},
		}
		return indexs
	}
	_, err := s.col.Indexes().CreateMany(context.Background(), indexs())
	s.log.Info().Msgf("the mongodb initialition is succeeds")

	domain.RegisterRPCServer(grpc.Get().Server(), s)
	return err
}
