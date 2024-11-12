package exception_test

import (
	"fmt"
	"testing"

	"github.com/kade-chen/library/exception"
)

func TestNewError(t *testing.T) {
	e := exception.NewNotFound("test %s", "ss")
	t.Log(e.ToJson())
	fmt.Println(e.Error())
}
