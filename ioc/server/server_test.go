package server_test

import (
	"context"
	"testing"

	"github.com/kade-chen/library/ioc/server"
)

func TestServer(t *testing.T) {
	server.Run(context.Background(), nil)
}
