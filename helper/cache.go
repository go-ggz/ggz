package helper

import (
	"strconv"
	"strings"
)

const sep = ":"

// GetCacheKey get cache key for data loader
func GetCacheKey(module string, id interface{}) string {
	var str string
	switch v := id.(type) {
	case int64:
		str = strconv.FormatInt(v, 10)
	case string:
		str = v
	}
	return module + sep + str
}

// GetCacheID get cache id for model id
func GetCacheID(key string) (interface{}, error) {
	strs := strings.Split(key, sep)

	return strconv.ParseInt(strs[1], 10, 64)
}
