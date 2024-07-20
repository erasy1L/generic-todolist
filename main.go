package main

import (
	"embed"
	"io/fs"
	"todo-list/internal/app"
	"todo-list/internal/server"
)

//go:embed api/docs/*
var dist embed.FS

// FS holds embedded swagger-ui files
var FS, _ = fs.Sub(dist, "api/docs")

func main() {
	server.FS = FS

	app.Run()
}
