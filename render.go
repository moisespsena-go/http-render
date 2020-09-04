package http_render

import (
	"context"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/moisespsena-go/httpu"
)

type contextKey struct{}

var ContextKey contextKey

type Render struct {
	Request                *http.Request
	Status                 int
	Data                   interface{}
	FileNames              []string
	AssetFactory           AssetFactory
	DirectoryIndexDisabled bool
	TemplateFuncMap        FuncMap
}

func New(opt ...Option) (r Render) {
	for _, opt := range opt {
		r = opt(r)
	}
	return
}

func (this Render) Option(opt ...Option) Render {
	for _, opt := range opt {
		this = opt(this)
	}
	return this
}

func (this Render) exec(dst http.ResponseWriter, fileName string) (err error) {
	var asset Asset
	if strings.HasSuffix(fileName, "/") {
		err = ErrIsDir
	} else {
		asset, err = this.AssetFactory.Factory(fileName)
	}
	if err != nil {
		if err == ErrIsDir && !this.DirectoryIndexDisabled {
			this.FileNames = []string{path.Join(fileName, "index")}
			return this.Render(dst)
		}
		return err
	}
	this.FileNames = []string{fileName}
	return asset.Execute(&this, dst)
}

func (this Render) Render(dst http.ResponseWriter) (err error) {
	if this.Data == nil {
		this.Data = map[string]interface{}{
			"Prefix":  httpu.Prefix(this.Request.Context()),
			"Request": this.Request,
		}
		var hasGet, hasPost bool

		if this.TemplateFuncMap != nil {
			if _, ok := this.TemplateFuncMap["GET"]; ok {
				hasGet = true
			}
			if _, ok := this.TemplateFuncMap["POST"]; ok {
				hasPost = true
			}
		}
		if !hasGet {
			var urlParam = this.Request.URL.Query()
			this.TemplateFuncMap.Update(map[string]interface{}{
				"GET": func(key string) string {
					return urlParam.Get(key)
				},
			})
		}
		if !hasPost {
			switch this.Request.Method {
			case http.MethodPost, http.MethodPut:
				this.TemplateFuncMap.Update(map[string]interface{}{
					"POST": this.Request.FormValue,
				})
			}
		}
	}

	for _, fileName := range this.FileNames {
		if err = this.exec(dst, fileName); err == nil {
			return
		}
		if os.IsNotExist(err) {
			if !strings.HasSuffix(fileName, "/") && !HasExt(fileName) {
				for _, ext := range []string{".tmpl", ".html"} {
					if err = this.exec(dst, fileName+ext); err == nil {
						return
					} else if !os.IsNotExist(err) {
						return
					}
				}
			}
		} else {
			return
		}
	}
	return os.ErrNotExist
}

func (this Render) RenderOrNotFound(w http.ResponseWriter) (err error) {
	if err = this.Render(w); os.IsNotExist(err) {
		http.NotFound(w, this.Request)
	}
	return
}

func (this Render) MustRenderOrNotFound(w http.ResponseWriter) {
	if err := this.Render(w); err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, this.Request)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (this Render) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), ContextKey, &this))
	this.Request = r
	this.FileNames = []string{r.URL.Path}
	this.MustRenderOrNotFound(w)
}
