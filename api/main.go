package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	c "cryptbox/server/common"
)

var rootPath = flag.String("root", ".", "Root folder to share")

type handler func(http.ResponseWriter, *http.Request) error

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request) error {
	path := filepath.Join(*rootPath, r.FormValue("path"))

	fmt.Println("Received list request for: ", path)

	dir, err := os.Open(path)
	if err != nil {
		return err
	}

	dirStat, err := dir.Stat()
	if err != nil {
		return err
	}

	if ! dirStat.IsDir() {
		return err
	}

	response := c.NewFileList(path)
	response.LoadFrom(dir)

	bytes, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := w.Write(bytes); err != nil {
		return err
	}

	return nil
}

func FileHandler(w http.ResponseWriter, r *http.Request) error {
	path := filepath.Join(*rootPath, r.FormValue("path"))

	fmt.Println("Received file request for: ", path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	if fileStat.IsDir() {
		return err
	}

	_, filename := filepath.Split(r.FormValue("path"))

	w.Header().Set("Content-Disposition", fmt.Sprint("attachment; filename=\"", filename, "\""))
	w.Header().Set("Content-Transfer-Encoding", "binary")

	http.ServeContent(w, r, filename, fileStat.ModTime(), file)
	return nil
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

	http.Handle("/list/", handler(ListHandler))
	http.Handle("/file/", handler(FileHandler))
	log.Fatal(http.ListenAndServe(":5555", nil))
}
