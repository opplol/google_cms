package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets6dc357a3c31dcc10799799ed15afd03f27f46286 = "<html>\n{{range .data_model}}\n{{.Title|safehtml}}\n{{.Body|safehtml}}\n{{end}}\n</html>"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"view"}, "/view": []string{"index.tmpl"}}, map[string]*assets.File{
	"/view/index.tmpl": &assets.File{
		Path:     "/view/index.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1627047460, 1627047460400000000),
		Data:     []byte(_Assets6dc357a3c31dcc10799799ed15afd03f27f46286),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1628784387, 1628784387931983800),
		Data:     nil,
	}, "/view": &assets.File{
		Path:     "/view",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1626531593, 1626531593830000000),
		Data:     nil,
	}}, "")
