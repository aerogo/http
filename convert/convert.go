package convert

// pow10
var pow10 = [...]int{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
}

// pow16
var pow16 = [...]int{
	1, 16, 256, 4096, 65536,
	1048576, 16777216, 268435456, 4294967296, 68719476736,
	1099511627776, 17592186044416, 281474976710656, 4503599627370496,
}

// ASCIIDecToInt converts an ascii encoded decimal integer to an int.
func ASCIIDecToInt(slice []byte) int {
	num := 0
	length := len(slice)

	for i := 0; i < length; i++ {
		num += (int(slice[i]) - 48) * pow10[length-i-1]
	}

	return num
}

// ASCIIHexToInt converts an ascii encoded hexadecimal integer to an int.
func ASCIIHexToInt(slice []byte) int {
	num := 0
	length := len(slice)

	for i := 0; i < length; i++ {
		code := int(slice[i])
		pos := length - i - 1

		if code >= 48 && code <= 57 {
			// Numbers
			num += (code - 48) * pow16[pos]
		} else if code >= 65 && code <= 70 {
			// Letters (uppercase)
			num += (code - 65 + 10) * pow16[pos]
		} else if code >= 97 && code <= 102 {
			// Letters (lowercase)
			num += (code - 97 + 10) * pow16[pos]
		}
	}

	return num
}
