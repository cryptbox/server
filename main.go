package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest"
	"net/http"
	"os"
	"path/filepath"
)

var rootPath = flag.String("Root Folder", ".", "Root folder to share")

type Object struct {
	IsDir bool `json:"isDir"`
	Name  string `json:"name"`
	Size  int64 `json:"size"`
}

func GetDirs(w *rest.ResponseWriter, req *rest.Request) {
	relPath, _ := base64.StdEncoding.DecodeString(req.PathParam("path"))
	path := filepath.Join(*rootPath, string(relPath))

	fmt.Println("Received request for: ", path)

	dir, err := os.Open(path)
	if err != nil {
		return
	}

	subDirs, err := dir.Readdir(100)
	results := make([]*Object, 0, len(subDirs))
	for _, subDir := range subDirs {
		results = append(results, &Object{subDir.IsDir(), subDir.Name(), subDir.Size()})
	}

	w.WriteJson(&results)
}

func main() {
	flag.Parse()

	if _, err := os.Stat(*rootPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Invalid root path: ", err)
		}
	}
	path, _ := filepath.Abs(*rootPath)
	fmt.Println("Serving: ", path)

	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		rest.Route{"GET", "/dirs/:path", GetDirs},
	)
	http.ListenAndServe(":5555", &handler)
}
