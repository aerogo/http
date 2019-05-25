package client_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/aerogo/http/client"
	"github.com/stretchr/testify/assert"
)

var urls = []string{
	// Popular stuff
	"https://google.com",
	"https://facebook.com",
	"https://twitter.com",
	"https://youtube.com",
	"https://naver.com",
	"https://cloudflare.com",

	// Queries
	"https://www.google.com/search?q=test",

	// Anime friends
	"https://notify.moe",
	"https://anilist.co",
	"https://myanimelist.net",
	"https://notify.moe/search/fate stay night",
	"https://notify.moe/search/fate%20stay%20night",
}

func testResponse(t *testing.T, response *client.Response, err error) {
	assert.NoError(t, err)
	assert.True(t, response.Ok())
	assert.NotZero(t, response.StatusCode())
	assert.Equal(t, response.RawLength(), len(response.Raw()))
	assert.Equal(t, response.RawLength(), len(response.RawString()))
	assert.NotEmpty(t, response.RawHeaders())
	assert.NotEmpty(t, response.RawHeadersString())

	redirect := response.HeaderString("Location")
	assert.True(t, len(response.String()) > 0 || redirect != "")

	buffer := bytes.Buffer{}
	n, err := response.WriteTo(&buffer)
	assert.NoError(t, err)
	assert.True(t, int(n) >= response.RawLength())
}

func TestClient(t *testing.T) {
	for _, url := range urls {
		println("URL", url)
		response, err := client.Get(url).End()
		testResponse(t, response, err)
		_, err = response.WriteTo(ioutil.Discard)
		assert.NoError(t, err)
	}
}

func TestClientNoGZip(t *testing.T) {
	for _, url := range urls {
		println("URL", url)
		response, err := client.Get(url).Header("Accept-Encoding", "identity").End()
		testResponse(t, response, err)
		_, err = response.WriteTo(ioutil.Discard)
		assert.NoError(t, err)
	}
}

func BenchmarkClient(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := client.Get(urls[0]).End()

		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkClientWithBody(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		response, _ := client.Get(urls[0]).End()

		if response.String() == "" {
			b.Error("Empty response")
		}
	}
}
