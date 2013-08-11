package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
	"path/filepath"
)

var rootPath = flag.String("Root Folder", ".", "Root folder to share")

type FileList struct {
	Path string `json:"path"`
	More bool   `json:"more"`
	Children []*File `json:"children"`
}

func NewFileList(path string) *FileList {
	response := new(FileList)
	response.Children = make([]*File, 0)

	return response
}

func (l *FileList) LoadFrom(dir *os.File) {
	subDirs, err := dir.Readdir(100)

	if err == io.EOF {
		l.More = false
	}

	for _, subDir := range subDirs {
		l.AddFile(&File{subDir.IsDir(), subDir.Name(), subDir.Size()})
	}
}

func (r *FileList) AddFile(o *File) {
	r.Children = append(r.Children, o)
}

type File struct {
	IsDir bool   `json:"isDir"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
}

func DirsHandler(w http.ResponseWriter, req *http.Request) {
	path := filepath.Join(*rootPath, req.FormValue("path"))

	fmt.Println("Received dir request for: ", path)

	dir, err := os.Open(path)
	if err != nil {
		return
	}

	dirStat, err := dir.Stat()
	if err != nil {
		return
	}

	response := NewFileList(path)
	if dirStat.IsDir() {
		response.LoadFrom(dir)
	}

	bytes, _ := json.Marshal(response)
	w.Write(bytes)
}

func FileHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprint("attachment; filename=\"", req.FormValue("path"), "\""))
	w.Header().Set("Content-Transfer-Encoding", "binary")

	fmt.Printf("%#v\n", w.Header())
	fmt.Println(*rootPath)

	path := filepath.Join(*rootPath, req.FormValue("path"))

	fmt.Println("Received file request for: ", path)

	dir, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to open file: ", path)
		return
	}

	dirStat, err := dir.Stat()
	if err != nil {
		fmt.Println("Unable to stat file: ", path)
		return
	}

	if dirStat.IsDir() {
		fmt.Println("Object is file: ", path)
		return
	}

	io.Copy(w, dir)
}

func main() {
	flag.Parse()

	if _, err := os.Stat(*rootPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Invalid root path: ", err)
		}
	}
	path, _ := filepath.Abs(*rootPath)
	rootPath = &path
	fmt.Println("Serving: ", path)

	http.HandleFunc("/dirs/", DirsHandler)
	http.HandleFunc("/file/", FileHandler)
	log.Fatal(http.ListenAndServe(":5555", nil))
}
