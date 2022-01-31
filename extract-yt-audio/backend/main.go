package main

import (
	"encoding/json"
	"extract-audio/videoservice"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var audioFile = ""

func GetURL(w http.ResponseWriter, r *http.Request) {

}

func Extract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	var videoReq videoservice.VideoReq
	json.NewDecoder(r.Body).Decode(&videoReq)

	videoInfo := videoservice.GetVideoInfo(videoReq.URL) // to get title of the video only
	ytURL := videoReq.URL
	audioTitle := videoInfo.Title
	fmt.Println("youtube video URL: ", ytURL)

	//command format:  youtube-dl --extract-audio --audio-format mp3 <link>
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "--extract-audio")
	cmdArgs = append(cmdArgs, "--audio-format")
	cmdArgs = append(cmdArgs, "mp3")
	cmdArgs = append(cmdArgs, "--output")
	cmdArgs = append(cmdArgs, audioTitle)
	cmdArgs = append(cmdArgs, ytURL)

	cmd := exec.Command("youtube-dl", cmdArgs...)
	stdout, _ := cmd.StdoutPipe()
	oneByte := make([]byte, 100)
	cmd.Start()

	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			break
		}
		r, _ := regexp.Compile("(100|(\\d{1,2}(\\.\\d+)*))%")

		downloadStatus := r.Find(oneByte)
		statusStr := string(downloadStatus)
		fmt.Fprintf(w, "%v\n", statusStr)
	}
	cmd.Wait()
	cmd.Process.Kill()
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

	// http.ServeContent(w, r, audioFile, time.Now(), bytes.NewReader(dat))
	http.ServeFile(w, r, audioFile)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/get-url", GetURL)
	router.HandleFunc("/extract-audio", Extract)

	port := ":8000"
	fmt.Printf("Server started @ %s\n", port)
	err := http.ListenAndServe(port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}
