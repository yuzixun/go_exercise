package utils

import (
	"bytes"
	"math"
	"strings"
)

var dict = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func decimalTo62(num int64) string {
	var builder strings.Builder

	for num != 0 {
		builder.WriteByte(dict[num%62])
		num = num / 62
	}

	return builder.String()
}

func decimalFrom62(base62 string) int64 {
	var num int64
	for i := 0; i < len(base62); i++ {
		num = num + int64(bytes.IndexByte(dict, base62[i]))*int64(math.Pow(62, float64(i)))
	}

	return num
}
