package zstring

import (
	"math/rand"
	"strings"
	"time"
)

func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func Empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func Rand(charset string, length int) string {
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rander.Intn(62)])
	}

	return sb.String()
}

var (
	_randDigits     = "0123456789"
	_randAlphabets  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_randAllCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func RandDigits(length int) string {
	return Rand(_randDigits, length)
}

func RandAlphabets(length int) string {
	return Rand(_randAlphabets, length)
}

func RandString(length int) string {
	return Rand(_randAllCharset, length)
}
