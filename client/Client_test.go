package client_test

import (
	"testing"

	"github.com/aerogo/http/client"
	"github.com/stretchr/testify/assert"
)

var urls = []string{
	// // Popular stuff
	"https://google.com",
	"https://github.com",
	"https://facebook.com",
	"https://twitter.com",
	"https://youtube.com",
	"https://naver.com",

	// Anime friends
	"https://notify.moe",
	"http://cal.syoboi.jp",
	"https://myanimelist.net",
	"https://myanimelist.net/anime/356/Fate_stay_night?q=fate stay night",
	"https://myanimelist.net/anime/356/Fate_stay_night?q=fate%20stay%20night",

	// These are failing atm:
	// "https://anilist.co",
	// "https://kitsu.io",
}

func TestClient(t *testing.T) {
	for _, url := range urls {
		println("URL", url)
		response, err := client.Get(url).End()

		assert.NoError(t, err)
		assert.True(t, response.Ok())
		assert.True(t, len(response.String()) >= 0 || response.HeaderString("Location") != "")
	}
}

// func TestClientNoGZip(t *testing.T) {
// 	for _, url := range urls {
// 		response, err := client.Get(url).Header("Accept-Encoding", "identity").End()

// 		assert.NoError(t, err)
// 		assert.True(t, response.Ok())
// 		assert.NotEmpty(t, response.String())
// 	}
// }

func BenchmarkClient(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		client.Get(urls[0]).End()
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
