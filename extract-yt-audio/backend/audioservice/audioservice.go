package audioservice

import (
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
)

func ConvertToAudio(ytURL string, audioTitle string, w http.ResponseWriter) {
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

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Internal error", 500)
		return
	}

	cmd.Start()

	for {
		if _, err := stdout.Read(oneByte); err != nil {
			break
		}
		r, _ := regexp.Compile("(100|(\\d{1,2}(\\.\\d+)*))%")

		downloadStatus := r.Find(oneByte)
		statusStr := string(downloadStatus)
		fmt.Fprintf(w, "data: %v\n\n", statusStr)
		flusher.Flush()
		fmt.Println(statusStr)
	}
	cmd.Wait()
	cmd.Process.Kill()
}
