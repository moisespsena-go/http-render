package ropt

import (
	"github.com/moisespsena-go/os-common"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/moisespsena-go/assetfs"

	"github.com/moisespsena-go/assetfs/assetfsapi"
	fsapi "github.com/moisespsena-go/assetfs/assetfsapi"
	http_render "github.com/moisespsena-go/http-render"
)

type FileSystemAssetFactory struct {
	FS fsapi.Interface
}

func (this FileSystemAssetFactory) Factory(name string) (exec http_render.Asset, err error) {
	var info fsapi.FileInfo
	if info, err = this.FS.AssetInfo(name); err != nil {
		if oscommon.IsNotFound(err) {
			err = os.ErrNotExist
		}
		return
	}
	if info.IsDir() {
		return nil, http_render.ErrIsDir
	}
	if strings.HasSuffix(name, ".tmpl") {
		return http_render.AssetFunc(func(r *http_render.Render, dst http.ResponseWriter) (err error) {
			var (
				tmpl   *template.Template
				reader io.ReadCloser
				b      []byte
			)

			if reader, err = info.Reader(); err != nil {
				return
			}

			if b, err = ioutil.ReadAll(reader); err != nil {
				return
			}

			if tmpl, err = template.New(info.RealPath()).Funcs(template.FuncMap(r.TemplateFuncMap)).Parse(string(b)); err != nil {
				return
			}

			if dst.Header().Get("Content-Type") == "" {
				dst.Header().Set("Content-Type",
					mime.FormatMediaType("text/html", map[string]string{"charset": "utf-8"}))
			}

			if r.Status != 0 {
				dst.WriteHeader(r.Status)
			}

			return tmpl.Execute(dst, r.Data)
		}), nil
	}
	var staticHandler = assetfs.NewStaticHandler(this.FS)
	return http_render.AssetFunc(func(r *http_render.Render, dst http.ResponseWriter) (err error) {
		if ext := http_render.Ext(name); ext != "" {
			if mt := mime.TypeByExtension(ext); mt != "" {
				dst.Header().Set("Content-Type", mt)
			}
		}
		r.Request.URL.Path = r.FileNames[0]
		staticHandler.ServeHTTP(dst, r.Request)
		return nil
	}), nil
}

func FS(fs assetfsapi.Interface) Option {
	return AssetFinder(&FileSystemAssetFactory{fs})
}
