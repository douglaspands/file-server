package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed templates/* static/*
var embeddedFiles embed.FS

// GetFileSystem retorna o sistema de arquivos para arquivos estáticos embutidos.
func GetStaticFileSystem() (http.FileSystem, error) {
	sub, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		return nil, err
	}
	return http.FS(sub), nil
}

// GetTemplatesFS retorna o sistema de arquivos contendo os templates.
func GetTemplatesFS() fs.FS {
	return embeddedFiles
}
