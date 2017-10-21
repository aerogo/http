package client_test

import (
	"testing"

	"github.com/aerogo/http/client"
	"github.com/parnurzeal/gorequest"
)

const url = "https://google.com"

func BenchmarkClient(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		client.Get(url).End()
	}
}

func BenchmarkClientWithBody(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := client.Get(url).End()
		resp.Body()
	}
}

func BenchmarkGoRequest(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		gorequest.New().Get(url).End()
	}
}
