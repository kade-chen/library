package time

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JwtTime1(end *jwt.NumericDate) int64 {
	nowPlus1h := time.Now().Add(time.Hour)

	var final time.Time
	if end == nil || nowPlus1h.Before(end.Time) {
		final = nowPlus1h
	} else {
		final = end.Time
	}
	fmt.Println("final", final)
	return final.Unix()
}

func JwtTime(end *jwt.NumericDate) int64 {
	now := time.Now()
	nowPlus1h := now.Add(time.Hour)

	var final time.Time
	if end == nil || nowPlus1h.Before(end.Time) {
		final = nowPlus1h
	} else {
		final = end.Time
	}
	fmt.Println("final", final)
	// 剩余秒数
	ttl := time.Until(final).Seconds()
	if ttl < 0 {
		return 0
	}

	return int64(ttl)
}
