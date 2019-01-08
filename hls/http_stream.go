package hls

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var streamRegexp = regexp.MustCompile(`^(.*)/([0-9]+)\.ts$`)
var localPath2 = "http://127.0.0.1:8888/";
type StreamHandler struct {
	root    string
	rootUri string
	encoder *Encoder
}

//func NewStreamHandler(root string, rootUri string) *StreamHandler {
//	return &StreamHandler{root, rootUri, NewEncoder("segments", 2)}
//}
func NewStreamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["fileId"];
	//tsId := vars["tsId"];

	log.Debugf("Stream request: %v", r.URL.Path)
	matches := streamRegexp.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		http.Error(w, "Wrong path format", 400)
		return
	}
	res := int64(720)
	segment, _ := strconv.ParseInt(matches[2], 0, 64)
	file := localPath2+fileId
	er := NewEncodingRequest(file, segment, res)
	NewEncoder("segments", 2).Encode(*er)

	w.Header()["Access-Control-Allow-Origin"] = []string{"*"}
	select {
	case data := <-er.data:
		w.Write(*data)
	case err := <-er.err:
		log.Errorf("Error encoding %v", err)
	case <-time.After(60 * time.Second):
		log.Errorf("Timeout encoding")
	}
}
