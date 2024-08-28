package zkvs

import (
	"strings"

	"github.com/yyliziqiu/zlib/zconv"
)

type KVS map[string]string

func (k KVS) ShouldGet(key string) (string, bool) {
	val, ok := k[key]
	if !ok {
		return "", false
	}
	return strings.TrimSpace(val), true
}

func (k KVS) ShouldBool(key string) (bool, bool) {
	val, ok := k.ShouldGet(key)
	if !ok {
		return false, false
	}
	return zconv.S2B(val), true
}

func (k KVS) ShouldInt(key string) (int, bool) {
	val, ok := k.ShouldGet(key)
	if !ok {
		return 0, false
	}
	return zconv.S2I(val), true
}

func (k KVS) ShouldInt64(key string) (int64, bool) {
	val, ok := k.ShouldGet(key)
	if !ok {
		return 0, false
	}
	return zconv.S2I64(val), true
}

func (k KVS) ShouldFloat64(key string) (float64, bool) {
	val, ok := k.ShouldGet(key)
	if !ok {
		return 0, false
	}
	return zconv.S2F64(val), true
}

func (k KVS) Get(key string) string {
	val, _ := k.ShouldGet(key)
	return val
}

func (k KVS) Bool(key string) bool {
	val, _ := k.ShouldBool(key)
	return val
}

func (k KVS) Int(key string) int {
	val, _ := k.ShouldInt(key)
	return val
}

func (k KVS) Int64(key string) int64 {
	val, _ := k.ShouldInt64(key)
	return val
}

func (k KVS) Float64(key string) float64 {
	val, _ := k.ShouldFloat64(key)
	return val
}

func (k KVS) Id() string {
	return k.Get("id")
}

func (k KVS) Name() string {
	return k.Get("name")
}
