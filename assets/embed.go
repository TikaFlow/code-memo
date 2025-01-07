package assets

import (
	"embed"
	"io/fs"
	"strings"
)

var (
	//go:embed public
	webStaticFiles embed.FS
	webStaticRoot  = "public"
	//go:embed static
	appStaticFiles embed.FS
	appStaticRoot  = "static"
	R              map[string][]byte
)

func init() {
	R = make(map[string][]byte)
	if err := fs.WalkDir(appStaticFiles, appStaticRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			content, err := appStaticFiles.ReadFile(path)
			if err != nil {
				return err
			}
			key := strings.TrimPrefix(path, appStaticRoot)
			R[key] = content
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func GetStatic() fs.FS {
	static, _ := fs.Sub(webStaticFiles, webStaticRoot)
	return static
}
