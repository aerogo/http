package client_test

import (
	"testing"

	"github.com/aerogo/http/client"
	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
)

const path = "https://blitzprog.org/"

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
		resp.String()
	}
}

func BenchmarkGoRequest(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		gorequest.New().Get(path).End()
	}
}
