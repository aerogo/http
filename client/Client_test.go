package client_test

import (
	"io/ioutil"
	"testing"

	"github.com/aerogo/http/client"
	"github.com/stretchr/testify/assert"
)

var urls = []string{
	// Popular stuff
	"https://google.com",
	"https://github.com",
	"https://facebook.com",
	"https://twitter.com",
	"https://youtube.com",
	"https://naver.com",
	"https://cloudflare.com",

	// Anime friends
	"https://notify.moe",
	"https://anilist.co",
	"https://myanimelist.net",
	"https://myanimelist.net/anime/356/Fate_stay_night?q=fate stay night",
	"https://myanimelist.net/anime/356/Fate_stay_night?q=fate%20stay%20night",
	"http://cal.syoboi.jp",

	// These are failing atm due to wrong status codes returned by the server:
	// "https://kitsu.io",
}

func testResponse(t *testing.T, response *client.Response, err error) {
	assert.NoError(t, err)
	assert.True(t, response.Ok())
	assert.True(t, len(response.String()) >= 0 || response.HeaderString("Location") != "")
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
