package utils

import (
	"math"
)

func EncodeRevertTimestamp(t int64) int64 {
	return math.MaxInt64 - t
}
func DecodeRevertTimestamp(rT int64) int64 {
	return math.MaxInt64 - rT
}
