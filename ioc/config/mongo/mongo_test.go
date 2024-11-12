package mongo_test

import (
	"os"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/mongo"
)

func TestGetClientGetter(t *testing.T) {
	m := mongo.Client()
	t.Log(m)
}

func TestGetClientGetter1(t *testing.T) {
	m := mongo.DB()
	t.Log(m)
}

func init() {
	os.Setenv("MONGO_ENDPOINTS", "127.0.0.1:55000")
	os.Setenv("MONGO_USERNAME", "admin")
	os.Setenv("MONGO_PASSWORD", "cys000522")
	os.Setenv("MONGO_DATABASE", "mcenter")
	os.Setenv("MONGO_AUTH_DB", "admin")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
