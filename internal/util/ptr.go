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
