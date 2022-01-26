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
	"github.com/r3labs/sse/v2"
)

var audioFile = ""

func ExtractAudio(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Reached ExtractAudio handler....")
	// w.Header().Set("Content-Type", "text/event-stream")
	// w.Header().Set("Cache-Control", "no-cache")
	// w.Header().Set("Connection", "Keep-alive")

	var videoReq videoservice.VideoReq
	json.NewDecoder(r.Body).Decode(&videoReq)

	sseServer := sse.New()
	sseServer.CreateStream("messages")

	fmt.Println("payload: ", videoReq)

	go func() {
		sseServer.Publish("messages", &sse.Event{
			Data: []byte("ping"),
		})
	}()
	sseServer.ServeHTTP(w, r)
	videoInfo := videoservice.GetVideoInfo(videoReq.URL) // to get title of the video only

	prgChan := make(chan string) // // audioTitleChan := make(chan string)
	var msg string

	go audioservice.ConvertToAudio(videoReq.URL, videoInfo.Title, prgChan)
	idx := 0
	for {
		idx += 1
		msg = <-prgChan
		fmt.Println("download status from chan: ", msg)
		if msg == "100%" {
			break
		}
		w.Write([]byte(msg))
	}
	json.NewEncoder(w).Encode(struct{ Msg string }{fmt.Sprint("download")})
}

func DownloadAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reached DownloadAudio handler....")
	fmt.Println("file to download: ", audioFile)

	_, err := ioutil.ReadFile(audioFile)
	if err != nil {
		log.Fatal(err)
	}

	// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "attachment; filename="+audioFile)
	w.Header().Set("Content-Transfer-Encoding", "binary")

	// http.ServeContent(w, r, audioFile, time.Now(), bytes.NewReader(dat))
	http.ServeFile(w, r, audioFile)
}

func Progress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "Keep-alive")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/extract-audio", ExtractAudio)
	router.HandleFunc("/download-audio", DownloadAudio)

	port := ":8000"
	fmt.Printf("Server started @ %s\n", port)
	err := http.ListenAndServe(port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}
