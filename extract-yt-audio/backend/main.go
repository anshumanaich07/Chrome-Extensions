package main

import (
	"encoding/json"
	"extract-audio/audioservice"
	"extract-audio/videoservice"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	audioURL  string
	audioFile string
)

func GetURL(w http.ResponseWriter, r *http.Request) {
	var videoReq videoservice.VideoReq
	json.NewDecoder(r.Body).Decode(&videoReq)
	audioURL = videoReq.URL

	fmt.Println("Received URL: ", audioURL)

	json.NewEncoder(w).Encode(struct{ Msg string }{fmt.Sprint("download")})
}

func Extract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	videoInfo := videoservice.GetVideoInfo(audioURL) // to get title of the video only
	audioFile = videoInfo.Title

	audioservice.ConvertToAudio(audioURL, videoInfo.Title, w)
}

func DownloadAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached DownloadAudio handler....")
	fmt.Println("file to download: ", audioFile)

	_, err := ioutil.ReadFile(audioFile)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+audioFile)
	w.Header().Set("Content-Transfer-Encoding", "binary")

	http.ServeFile(w, r, audioFile)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/get-url", GetURL)
	router.HandleFunc("/extract-audio", Extract)
	router.HandleFunc("/download-audio", DownloadAudio)

	port := ":8000"
	fmt.Printf("Server started @ %s\n", port)
	err := http.ListenAndServe(port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}
