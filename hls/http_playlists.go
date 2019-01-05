package hls

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)
type PlaylistHandler struct {
	root         string
	rootUri      string
	segmentsPath string
}

func NewPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["fileId"];
	fmt.Println(fileId)
	log.Debugf("Playlist request: %v", r.URL.Path)

	vinfo, err := GetVideoInformation(localPath+fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	duration := vinfo.Duration
	baseurl := fmt.Sprintf("http://%v", r.Host)


	w.Header()["Content-Type"] = []string{"application/vnd.apple.mpegurl"}
	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}

	fmt.Fprint(w, "#EXTM3U\n")
	fmt.Fprint(w, "#EXT-X-VERSION:3\n")
	fmt.Fprint(w, "#EXT-X-MEDIA-SEQUENCE:0\n")
	fmt.Fprint(w, "#EXT-X-ALLOW-CACHE:YES\n")
	fmt.Fprint(w, "#EXT-X-TARGETDURATION:"+fmt.Sprintf("%.f", hlsSegmentLength)+"\n")
	fmt.Fprint(w, "#EXT-X-PLAYLIST-TYPE:VOD\n")

	leftover := duration
	segmentIndex := 0
	for leftover > 0 {
		if leftover > hlsSegmentLength {
			fmt.Fprintf(w, "#EXTINF: %f,\n", hlsSegmentLength)
		} else {
			fmt.Fprintf(w, "#EXTINF: %f,\n", leftover)
		}
		fmt.Fprintf(w, baseurl+segmentsPath+"%v/%v.ts\n", fileId, segmentIndex)
		segmentIndex++
		leftover = leftover - hlsSegmentLength
	}
	fmt.Fprint(w, "#EXT-X-ENDLIST\n")
}
