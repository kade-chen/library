package redis_test

import (
	"context"
	"log"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/redis"
)

func TestRedisClient(t *testing.T) {

	if err := redis.Client().Ping(context.TODO()).Err(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}
	log.Println("✅ Redis connected")

}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "etc/application.toml"
	ioc.DevelopmentSetup(req)
}
