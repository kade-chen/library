package hash_test

import (
	"testing"

	"github.com/kade-chen/library/tools/hash"
)

func TestFNV(t *testing.T) {
	s := hash.FnvHash("123456", "ssss")
	t.Log(s)
}
