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

func (k KVS) Get(key string) string {
	return k.StringN(key, "")
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
