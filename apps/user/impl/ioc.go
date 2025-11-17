package impl

import (
	"context"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kade-chen/google-billing-console/apps/user"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/grpc"
	logs "github.com/kade-chen/library/ioc/config/log"
	mongodb "github.com/kade-chen/library/ioc/config/mongo"
)

var _ user.Service = &service{}

func init() {
	ioc.Controller().Registry(&service{})
}

type service struct {
	ioc.ObjectImpl
	col *mongo.Collection
	log *zerolog.Logger

	user.UnimplementedRPCServer
}

func (i *service) Name() string {
	return user.AppName
}

func (s *service) Init() error {
	s.col = mongodb.DB().Collection(user.AppName)

	//create index
	indexs := func() []mongo.IndexModel {
		var Unique bool = true
		indexs := []mongo.IndexModel{
			{Keys: bson.D{{Key: "create_at", Value: 1}}},
			{
				Keys: bson.D{
					{Key: "domain", Value: -1},
					{Key: "username", Value: -1},
				},
				Options: &options.IndexOptions{Unique: &Unique},
			},
		}
		return indexs
	}

	// create index
	// indexs := []mongo.IndexModel{
	// 	{Keys: bson.D{{Key: "create_at", Value: 1}}},
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "domain", Value: -1},
	// 			{Key: "username", Value: -1},
	// 		},
	// 		Options: &options.IndexOptions{Unique: &myBool1},
	// 	},
	// }

	_, err := s.col.Indexes().CreateMany(context.Background(), indexs())
	if err != nil {
		return err
	}
	s.log = logs.Sub(user.AppName)
	user.RegisterRPCServer(grpc.Get().Server(), s)

	// conn, _ := grpc1.Dial("localhost:1234", grpc1.WithTransportCredentials(insecure.NewCredentials()))
	// v := user.NewRPCClient(conn)
	// v.
	return nil
}
