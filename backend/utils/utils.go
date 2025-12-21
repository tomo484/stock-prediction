package utils

import (
	"strconv"
	"strings"
)

// "5.4201%" -> 5.4201 (float64)
func ParsePercentage(s string) float64 {
	s = strings.ReplaceAll(s, "%", "")
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

// "50.57" -> 50.57 (float64)
func ParseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

// "12345" -> 12345 (int64)
func ParseInt(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

// GetStringFromMap map[string]interface{}からstringを安全に取得する
// キーが存在しない、または型が一致しない場合は空文字列を返す
func GetStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetStringFromMapWithDefault map[string]interface{}からstringを取得し、存在しない場合はデフォルト値を返す
func GetStringFromMapWithDefault(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

