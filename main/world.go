package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// 1、实现读取文件的 Handler
	fileHandler := http.FileServer(http.Dir("./video"))

	// 2、注册handler
	http.Handle("/video/", fileHandler)

	//启动 web 服务
	http.ListenAndServe(":8090", nil)
}

// 1、上传视频文件接口的业务逻辑
func uploadHandler(w http.ResponseWriter, r *http.Request)  {
	//1、限制客户端上传视频文件的大小
	r.Body = http.MaxBytesReader(w, r.Body, 10 *1024 * 1024)
	err := r.ParseMultipartForm(10 * 1024 *1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//2、获取上传的文件
	file, fileHeader, err := r.FormFile("uploadFile")

	// 3、检查文件类型
	ret := strings.HasSuffix(fileHeader.Filename, ".mp4")
	if ret == false {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 4、获取随机名称
	// 文件名和当前时间的和
	md5Byte := md5.Sum([]byte(fileHeader.Filename + time.Now().String()))
	// 转换为 16 进制
	md5str := fmt.Sprintf("%x", md5Byte)
	newFileName := md5str + ".mp4"

	// 5、写入文件
	dst, err := os.Create("./video" + newFileName)
	defer dst.Close()
	if err != nll{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if _, err := io.Copy(dst, file); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// 获取视频文件列表
func getFileListHandler(w http.ResponseWriter, r *http.Request)  {
	// 通过filepath。Golb() 函数获取指定目录下的所以的文件，返回与非数组。数组的元素时每一个文件的路径
	files,_ := filepath.Glob("video/*")
	// 遍历数组，将每一个文件名改成 http 请求的 url 形式
	// 声明一个返回值的类型，是一个数组的切片
	var ret []string
	for _, file := range files {
		//r.Host 获取域名端口号， filepath.Base() 获取文件名
		ret = append(ret, "http://" + r.Host + "/video/" + filepath.Base(file))
	}
	// 然后将这个切牌你转换为 json 格式返回
	retJson, _ := json.Marshal(ret)
	// 将Json 写入到响应中进行返回
	w.Write(retJson)
	return
}
