package client_test

import (
	"testing"

	"github.com/aerogo/http/client"
	"github.com/stretchr/testify/assert"
)

const path = "https://github.com"

func TestClient(t *testing.T) {
	resp, err := client.Get(path).End()

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
	assert.NotEmpty(t, resp.String())
}

func BenchmarkClient(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		client.Get(path).End()
	}
}

func BenchmarkClientWithBody(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := client.Get(path).End()

		if resp.String() == "" {
			b.Error("Empty response")
		}
	}
}
