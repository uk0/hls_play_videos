package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/uk0/Cloud_Disk/hls"
	"github.com/uk0/Cloud_Disk/tools"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var localPath = "http://127.0.0.1:8888/";

func main() {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("app")))
	// 使用中
	r.HandleFunc("/uploadFile", uploadFile)
	// 使用中
	r.HandleFunc("/getFiles", getFiles)
	// 使用中
	r.HandleFunc("/getdirs", getDirs)
	// 使用中
	r.HandleFunc("/mkdir", mkdir)
	// 使用中
	r.HandleFunc("/download/{fid}", download)
	// 废弃
	r.HandleFunc("/list", getList)
	// 使用中
	r.HandleFunc("/playlist/{fileId}", hls.NewPlaylistHandler)
	// 使用中
	r.HandleFunc("/frame", hls.NewFrameHandler)
	// 使用中
	r.HandleFunc("/segments/{fileId}/{tsId}", hls.NewStreamHandler)
	http.Handle("/", r)
	// 使用中（看似多余，实际是为了区分 ）
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("app"))))
	//tools.InitSqlDB();
	http.ListenAndServe(":8001", nil)
}

func download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fid := vars["fid"];
	w.Write([]byte(fid))
}


func getFiles(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryForm["fid"]) == 0 {
		w.Write([]byte("path?fid="))
	}
	fid := queryForm["fid"][0];
	w.Write(tools.GetFileByFid(fid))
}

func mkdir(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryForm["dir"]) == 0 {
		w.Write([]byte("path?dir=/"))
	}
	path := queryForm["dir"][0];
	w.Write([]byte(tools.Mkdir(path)))
}

func getDirs(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryForm["dir"]) == 0 {
		w.Write([]byte("path?dir=/"))
	}
	path := queryForm["dir"][0];
	log.Println(path)
	e, err := json.Marshal(tools.GetDirs(path))
	if err != nil {
		log.Println("ERROR")
	}
	w.Write(e)
}

func getList(w http.ResponseWriter, r *http.Request) {
	var tempData = tools.QueryList();
	e, err := json.Marshal(tempData)
	if err != nil {
		log.Println("ERROR")
	}
	w.Write(e)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryForm["dir"]) == 0 {
		w.Write([]byte("path?dir=/"))
	}
	path := queryForm["dir"][0];
	// 根据字段名获取表单文件
	formFile, header, err := r.FormFile("uploadfile")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()

	//var blok = tools.GetBlokInfo();
	//var Fid = tools.GetFid(blok)

	var saveFileNmae = strings.Replace(string(header.Filename)," ","_",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"'","_",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"\"","_",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"）",")",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"（","(",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"，","_",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"-","_",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"“","_",-1)
	    saveFileNmae = strings.Replace(saveFileNmae,"：","_",-1)

	log.Println("File For Get Header.Filename: " + saveFileNmae)
	url := tools.PublicUrl + path + "/" + saveFileNmae
	log.Println("url : " + url)
	go tools.PutFile(url, formFile, saveFileNmae)
	//log.Printf("fileSize: %s\n", fileSize)
	//tools.SaveInfo(Fid, string(time.Now().Format("2006-01-02 15:04:05")), header.Filename, fileSize)
	log.Printf("FileName: %s\n", saveFileNmae)
	w.Write([]byte("success"));
}
