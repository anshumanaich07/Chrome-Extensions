package main

import (
	"bytes"
	"encoding/json"
	"extract-audio/audioservice"
	"extract-audio/videoservice"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var audioFile = ""

func ExtractAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached ExtractAudio handler....")

	var videoResponse videoservice.VideoResponse
	json.NewDecoder(r.Body).Decode(&videoResponse)

	videoInfo := videoservice.GetVideoInfo(videoResponse.URL)
	audioFile = audioservice.ConvertToAudio(videoResponse.URL, videoInfo.Title)

	json.NewEncoder(w).Encode(struct{ Msg string }{fmt.Sprint("download")})
}

func DownloadAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached DownloadAudio handler....")
	fmt.Println("file to download: ", audioFile)

	dat, err := ioutil.ReadFile(audioFile)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "attachment; filename="+audioFile)
	w.Header().Set("Content-Transfer-Encoding", "binary")

	http.ServeContent(w, r, audioFile, time.Now(), bytes.NewReader(dat))
	// http.ServeFile(w, r, audioFile)

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/extract-audio", ExtractAudio)
	r.HandleFunc("/download-audio", DownloadAudio)

	port := ":8000"
	fmt.Printf("Server started @ %s\n", port)
	err := http.ListenAndServe(port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
	if err != nil {
		log.Fatal(err)
	}

}
