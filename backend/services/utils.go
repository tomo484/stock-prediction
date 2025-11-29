package services

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