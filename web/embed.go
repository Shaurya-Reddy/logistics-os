package web

import (
	"embed"
	"io/fs"
)

// Build the frontend before compiling Go. Build output is deliberately untracked.
//
//go:embed dist/index.html dist/assets
var assets embed.FS

func Assets() fs.FS {
	root, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}
	return root
}
