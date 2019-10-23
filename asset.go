package http_render

import (
	"errors"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var ErrIsDir = errors.New("is directory")

type Asset interface {
	Execute(r *Render, dst http.ResponseWriter) (err error)
}

type AssetFunc func(r *Render, dst http.ResponseWriter) (err error)

func (this AssetFunc) Execute(r *Render, dst http.ResponseWriter) (err error) {
	return this(r, dst)
}

type AssetFactory interface {
	Factory(name string) (exec Asset, err error)
}

type DirectoryAssetFactory struct {
	Dir string
}

func (this DirectoryAssetFactory) GetDir() string {
	return this.Dir
}

func (this DirectoryAssetFactory) Factory(name string) (exec Asset, err error) {
	pth := filepath.Join(this.Dir, filepath.FromSlash(path.Clean(name)))
	if !strings.HasPrefix(pth, this.Dir) {
		return nil, fmt.Errorf("bad asset name")
	}
	var s os.FileInfo
	if s, err = os.Stat(pth); err != nil {
		return
	}
	if s.IsDir() {
		return nil, ErrIsDir
	}
	if strings.HasSuffix(pth, ".tmpl") {
		return AssetFunc(func(r *Render, dst http.ResponseWriter) (err error) {
			var tmpl *template.Template
			if tmpl, err = template.ParseFiles(pth); err != nil {
				return
			}
			if r.TemplateFuncMap != nil {
				tmpl = tmpl.Funcs(template.FuncMap(r.TemplateFuncMap))
			}

			dst.Header().Set("Content-Type", "text/html")

			if r.Status != 0 {
				dst.WriteHeader(r.Status)
			}

			return tmpl.Execute(dst, r.Data)
		}), nil
	}
	return AssetFunc(func(r *Render, dst http.ResponseWriter) (err error) {
		if ext := Ext(name); ext != "" {
			if mt := mime.TypeByExtension(ext); mt != "" {
				dst.Header().Set("Content-Type", mt)
			}
		}
		http.ServeFile(dst, r.Request, pth)
		return nil
	}), nil
}

func Ext(name string) (ext string) {
	if pos := strings.LastIndexByte(name, '.'); pos > 0 {
		ext = name[pos:]
		if !strings.HasSuffix(name, ext) {
			return
		}
	}
	return
}

func HasExt(name string) bool {
	return Ext(name) != ""
}
