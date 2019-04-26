package client_test

import (
	"testing"

	"github.com/aerogo/http/client"
	"github.com/stretchr/testify/assert"
)

const path = "http://localhost:4000"

func TestClient(t *testing.T) {
	response, err := client.Get(path).End()

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode())
	assert.NotEmpty(t, response.String())
}

func TestClientNoGZip(t *testing.T) {
	response, err := client.Get(path).Header("Accept-Encoding", "identity").End()

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode())
	assert.NotEmpty(t, response.String())
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
		response, _ := client.Get(path).End()

		if response.String() == "" {
			b.Error("Empty response")
		}
	}
}
