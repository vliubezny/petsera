package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var embededFiles embed.FS

func LoadFileSystem() (http.FileSystem, error) {
	root, err := fs.Sub(embededFiles, "dist")
	if err != nil {
		return nil, err
	}

	return http.FS(root), nil
}
