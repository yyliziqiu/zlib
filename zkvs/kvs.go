package zkvs

import (
	"strings"

	"github.com/yyliziqiu/zlib/zconv"
)

type KVS map[string]string

// 1

func (k KVS) String(key string) (string, bool) {
	if val, ok := k[key]; ok {
		return strings.TrimSpace(val), true
	}
	return "", false
}

func (k KVS) Bool(key string) (bool, bool) {
	if val, ok := k.String(key); ok {
		return zconv.S2B(val), true
	}
	return false, false
}

func (k KVS) Int(key string) (int, bool) {
	if val, ok := k.String(key); ok {
		return zconv.S2I(val), true
	}
	return 0, false
}

func (k KVS) Int64(key string) (int64, bool) {
	if val, ok := k.String(key); ok {
		return zconv.S2I64(val), true
	}
	return 0, false
}

func (k KVS) Float64(key string) (float64, bool) {
	if val, ok := k.String(key); ok {
		return zconv.S2F64(val), true
	}
	return 0, false
}

// 2

func (k KVS) StringN(key string, def string) string {
	if val, ok := k.String(key); ok {
		return val
	}
	return def
}

func (k KVS) BoolN(key string, def bool) bool {
	if val, ok := k.Bool(key); ok {
		return val
	}
	return def
}

func (k KVS) IntN(key string, def int) int {
	if val, ok := k.Int(key); ok {
		return val
	}
	return def
}

func (k KVS) Int64N(key string, def int64) int64 {
	if val, ok := k.Int64(key); ok {
		return val
	}
	return def
}

func (k KVS) Float64N(key string, def float64) float64 {
	if val, ok := k.Float64(key); ok {
		return val
	}
	return def
}

// 3

var lower = strings.ToLower

func (k KVS) StringV(key string) (string, bool) {
	return k.String(lower(key))
}

func (k KVS) BoolV(key string) (bool, bool) {
	return k.Bool(lower(key))
}

func (k KVS) IntV(key string) (int, bool) {
	return k.Int(lower(key))
}

func (k KVS) Int64V(key string) (int64, bool) {
	return k.Int64(lower(key))
}

func (k KVS) Float64V(key string) (float64, bool) {
	return k.Float64(lower(key))
}

func (k KVS) StringNV(key string, def string) string {
	return k.StringN(lower(key), def)
}

func (k KVS) BoolNV(key string, def bool) bool {
	return k.BoolN(lower(key), def)
}

func (k KVS) IntNV(key string, def int) int {
	return k.IntN(lower(key), def)
}

func (k KVS) Int64NV(key string, def int64) int64 {
	return k.Int64N(lower(key), def)
}

func (k KVS) Float64NV(key string, def float64) float64 {
	return k.Float64N(lower(key), def)
}

// 4

func (k KVS) Get(key string) string {
	return k.StringN(key, "")
}

func (k KVS) GetV(key string) string {
	return k.Get(lower(key))
}

func (k KVS) Id() string {
	return k.StringN("id", "")
}

func (k KVS) Name() string {
	return k.StringN("name", "")
}

func (k KVS) Enabled() bool {
	return k.BoolN("enabled", false)
}

func (k KVS) Disabled() bool {
	return k.BoolN("disabled", false)
}
