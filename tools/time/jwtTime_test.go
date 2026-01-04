package time_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	tools "github.com/kade-chen/google-billing-console/tools/time"
)

// ：end < now+1h → 返回 end
func TestJwtTime_EndEarlierThan1h(t *testing.T) {
	end := jwt.NewNumericDate(time.Now().Add(30 * time.Minute))
	fmt.Println("end", end)
	got := tools.JwtTime(end)
	fmt.Println("got", got)
}

// end > now+1h → 返回 now+1h
func TestJwtTime_EndLaterThan1h(t *testing.T) {
	start := time.Now()
	end := jwt.NewNumericDate(start.Add(2 * time.Hour))
	fmt.Println(end)
	got := tools.JwtTime(end)
	fmt.Println("got", got)
}

// end == nil → 返回 now+1h
func TestJwtTime_EndNil(t *testing.T) {

	got := tools.JwtTime(nil)
	fmt.Println("got", got)
}
