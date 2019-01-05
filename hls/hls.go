package hls

import (
	"os"
	"path/filepath"
)

var HomeDir = ".hls_cache"
var FFProbePath = "ffprobe"
var FFMPEGPath = "ffmpeg"
var segmentsPath = "/segments/"
var localPath = "http://127.0.0.1:8888/";

const cacheDirName = "cache"
const hlsSegmentLength = 5.0 // Seconds

func ClearCache() error {
	var cacheDir = filepath.Join(HomeDir, cacheDirName)
	return os.RemoveAll(cacheDir)
}
