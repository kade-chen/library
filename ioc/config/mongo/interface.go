package mongo

import (
	"github.com/kade-chen/library/ioc"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	AppName = "mongodb"
)

func DB() *mongo.Database {
	return ioc.Config().Get(AppName).(*mongoDB).GetDB()
}

func Client() *mongo.Client {
	return ioc.Config().Get(AppName).(*mongoDB).Client()
}

// 需要依赖 mongodb, 通过Ioc的依赖注入 来通过配置文件自动注入该对象
// //i.col = ioc_mongo.Client().Database(resource.AppName).Collection(resource.AppName)
// s, _ := ioc_mongo.Client().StartSession()
// s.StartTransaction()
// s.Client().Database(resource.AppName).Collection(resource.AppName)
// i.col = ioc_mongo.DB().Collection(resource.AppName)
