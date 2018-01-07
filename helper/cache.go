package helper

import (
	"strconv"
	"strings"
)

const sep = ":"

// GetCacheKey get cache key for data loader
func GetCacheKey(module string, id interface{}) string {
	var str string
	switch id.(type) {
	case string:
		str = id.(string)
	case int64:
		str = strconv.FormatInt(id.(int64), 10)
	}
	return module + sep + str
}

// GetCacheID get cache id for model id
func GetCacheID(key string) (interface{}, error) {
	strs := strings.Split(key, sep)
	switch strs[0] {
	case "user":
		return strconv.ParseInt(strs[1], 10, 64)
	default:
		return nil, nil
	}
}
