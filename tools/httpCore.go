package tools

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

var MasterUrl = "http://localhost:9333";


var PublicUrl = "http://localhost:8888";

func GetFid(t string) string {
	var x Data
	err := json.Unmarshal([]byte(t), &x)
	if err != nil {
		panic(err)
	}
	return x.Fid
}

func GetPaths(t string) Paths {
	var x Paths
	err := json.Unmarshal([]byte(t), &x)
	if err != nil {
		panic(err)
	}
	return x
}

func GetFileSize(t string) int32 {
	var x Resp
	err := json.Unmarshal([]byte(t), &x)
	if err != nil {
		panic(err)
	}
	return x.Size
}

type Paths struct {
	Path string `json:"Path"`
	Entries []Entries `json:"Entries"`
}

type Entries struct {
	FullPath string `json:"FullPath"`
	Mtime string `json:"Mtime"`
	Crtime string `json:"Crtime"`
	Mode int64 `json:"Mode"`
	Chunks []interface{} `json:"chunks"`
}



type Resp struct {
	FileName string `json:"name"`
	Size     int32  `json:"size"`
}

type Data struct {
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int32  `json:"count"`
}

// 获得文件块位置
func GetBlokInfo() string {
	url := MasterUrl + "/dir/assign?dataCenter=DefaultDataCenter"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body);
}

func Mkdir(dir string)string{
	url := "http://localhost:8888"+dir+"/root"

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body);
}

// {"mtime":1538882343,"size":110351,"name":"12b40.jpg","path":"12b40.jpg","is_dir":false,"is_deleteable":false,"is_readable":true,"is_writable":true,"is_executable":false}
//获取文件
func GetFile(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)

}

//获取文件
func GetFileByFid(Fid string) []byte {
	url := PublicUrl + "/" + Fid
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body;
}

func GetDirs(path string) Paths {
	url := "http://127.0.0.1:8888"+path+"?&pretty=y"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return  GetPaths(string(body))
}

// 提交文件

func PutFile(url string, re io.Reader, fileName string) int32 {
	// 创建表单文件
	// CreateFormFile 用来创建表单，第一个参数是字段名，第二个参数是文件名
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("uploadfile", fileName)
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
	}

	_, err = io.Copy(formFile, re)
	if err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
	}

	// 发送表单
	contentType := writer.FormDataContentType()
	writer.Close() // 发送之前必须调用Close()以写入结尾行
	resp, err := http.Post(url, contentType, buf)
	resp.Header.Add("accept", "application/json")
	resp.Header.Add("content-type", "application/json")
	if err != nil {
		log.Fatalf("Post failed: %s\n", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return GetFileSize(string(body));
}
