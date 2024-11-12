package enum_test

import (
	"testing"

	"github.com/kade-chen/library/cmd/generate/enum"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	should := assert.New(t)
	code, err := enum.G.Generate("../../../pb/example/test.pb.go")
	t.Log(string(code))
	should.NoError(err)
}
