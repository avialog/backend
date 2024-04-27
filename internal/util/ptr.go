package util

import "time"

func String(s string) *string {
	return &s
}

func Duration(t time.Duration) *time.Duration {
	return &t
}

func Uint(u uint) *uint {
	return &u
}

func Int64(i int64) *int64 {
	return &i
}
