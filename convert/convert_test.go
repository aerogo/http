package convert_test

import (
	"testing"

	"github.com/aerogo/http/convert"
	"github.com/stretchr/testify/assert"
)

func TestDecToInt(t *testing.T) {
	assert.Equal(t, 0, convert.ASCIIDecToInt([]byte("")))
	assert.Equal(t, 0, convert.ASCIIDecToInt([]byte("0")))
	assert.Equal(t, 1, convert.ASCIIDecToInt([]byte("1")))
	assert.Equal(t, 10, convert.ASCIIDecToInt([]byte("10")))
	assert.Equal(t, 100, convert.ASCIIDecToInt([]byte("100")))
	assert.Equal(t, 123456789, convert.ASCIIDecToInt([]byte("123456789")))
	assert.Equal(t, 0, convert.ASCIIDecToInt([]byte("ZZZ")))
}

func TestHexToInt(t *testing.T) {
	assert.Equal(t, 0x0, convert.ASCIIHexToInt([]byte("")))
	assert.Equal(t, 0x0, convert.ASCIIHexToInt([]byte("0")))
	assert.Equal(t, 0x1, convert.ASCIIHexToInt([]byte("1")))
	assert.Equal(t, 0xA, convert.ASCIIHexToInt([]byte("A")))
	assert.Equal(t, 0x10, convert.ASCIIHexToInt([]byte("10")))
	assert.Equal(t, 0xAFFE, convert.ASCIIHexToInt([]byte("AFFE")))
	assert.Equal(t, 0xAFFE, convert.ASCIIHexToInt([]byte("Affe")))
	assert.Equal(t, 0xBADFACE, convert.ASCIIHexToInt([]byte("BADFACE")))
	assert.Equal(t, 0xBADFACE, convert.ASCIIHexToInt([]byte("BadFace")))
	assert.Equal(t, 0, convert.ASCIIHexToInt([]byte("ZZZ")))
}
