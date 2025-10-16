package format_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kade-chen/library/tools/format"
)

func TestFormatJSON(t *testing.T) {
	var data struct {
		IntVal    int
		FloatVal  float64
		StringVal string
		BoolVal   bool
		SliceVal  []int
		MapVal    map[string]float64
		Nested    struct {
			X int
			Y float64
		}
		PtrVal   *string
		TimeVal  time.Time
		NilSlice []string
		NilMap   map[string]int
	}

	// 填充数据
	data.IntVal = 42
	data.FloatVal = 3.14159265358979323846141592653589793141592653589793141592653589793141592653589793141592653589793141592653589793
	data.StringVal = "Hello"
	data.BoolVal = true
	data.SliceVal = []int{1, 2, 3}
	data.MapVal = map[string]float64{"pi": 3.1415926535, "e": 2.7182818284}
	data.Nested.X = 100
	data.Nested.Y = 1.23456789
	str := "pointer string"
	data.PtrVal = &str
	data.TimeVal = time.Date(2025, 10, 16, 10, 30, 0, 0, time.UTC)

	fmt.Println("=== ToJSONV2 ===")
	fmt.Println(format.ToJSONV2(&data))

	fmt.Println("\n=== ToJSON (v1) ===")
	fmt.Println(format.ToJSON(&data))
}
