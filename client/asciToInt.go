package client

import (
	"math"
)

func asciiToInt(slice []byte) int {
	num := 0

	for i := 0; i < len(slice); i++ {
		num += (int(slice[i]) - 48) * int(math.Pow10(len(slice)-i-1))
	}

	return num
}
