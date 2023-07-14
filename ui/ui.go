package ui

import (
	"embed"
	"io/fs"
)

//go:embed dist
var dist embed.FS
var FS, Error = fs.Sub(dist, "dist")
