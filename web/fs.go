package web

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distEmbed embed.FS

func GetDistFS() (fs.FS, error) {
	return fs.Sub(distEmbed, "dist")
}
