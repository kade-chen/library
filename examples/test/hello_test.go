package test_test

import (
	"testing"

	hello1 "github.com/kade-chen/library/test"
	"github.com/kade-chen/library/test2"
)

func TestHello(t *testing.T) {
	hello1.Hello111()
	test2.Hello111()
}
