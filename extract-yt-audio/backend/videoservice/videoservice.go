package videoservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type VideoReq struct {
	URL string `json:"videoURL"`
}

const (
	videoInfoURL string = "https://noembed.com/embed?url=" // 3rd party API URL to get YT video info
)

type VideoInfo struct {
	URL             string `json:"url"`
	Version         string `json:"version"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	ProviderName    string `json:"provider_name"`
	Width           int    `json:"width"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	Title           string `json:"title"`
	ThumbnailURL    string `json:"thumbnail_url"`
	AuthorName      string `json:"author_name"`
	AuthorURL       string `json:"author_url"`
	HTML            string `json:"html"`
	Height          int    `json:"height"`
	ProviderURL     string `json:"provider_url"`
	Type            string `json:"type"`
}

func GetVideoInfo(ytURL string) (videoInfo VideoInfo) {
	x := videoInfoURL + ytURL
	d, err := http.Get(x)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(d.Body)
	json.Unmarshal([]byte(data), &videoInfo)

	t := []string{videoInfo.Title}

	videoInfo.Title = strings.Join(t, "") + ".mp3"

	return
}
