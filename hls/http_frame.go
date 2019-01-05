package hls

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)



func NewFrameHandler(w http.ResponseWriter, r *http.Request) {
	fileId := r.URL.Query().Get("fileId")
	t := r.URL.Query().Get("t")
	fmt.Println(t)
	time := 30
	if tint, err := strconv.Atoi(t); err == nil {
		time = tint
	}
	path :=  localPath + fileId;
	args := []string{
		"-timelimit", "15",
		"-loglevel", "error",
		"-ss", fmt.Sprintf("%v.0", time),
		"-i", path,
		"-vf", "scale=320:-1",
		"-frames:v", "1",
		"-f", "image2",
		"-",
	}
	if err := NewHttpCommandHandler(2, "frames").ServeCommand(FFMPEGPath, args, calculateCommandHash(FFMPEGPath, args), w); err != nil {
		log.Errorf("Problem serving screenshot: %v", err)
	}
}
