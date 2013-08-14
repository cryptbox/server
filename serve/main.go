package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var rootPath = flag.String("root", ".", "Root folder to share")

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

	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(*rootPath))))
}
