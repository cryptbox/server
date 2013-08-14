package main

import (
	"io"
	"os"
)

type FileList struct {
	Path string `json:"path"`
	More bool   `json:"more"`
	Children []*File `json:"children"`
}

type File struct {
	IsDir bool   `json:"isDir"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
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

