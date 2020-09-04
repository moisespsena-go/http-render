package rrhandler

import (
	"net/http"
	"os"
	"path"
	"strings"

	http_render "github.com/moisespsena-go/http-render"
	"github.com/moisespsena-go/http-render/ropt"
)

type RequestRenderHandler struct {
	RootDir                  string
	RenderOrNotFoundDisabled bool
	Render                   http_render.Render
	RenderFunc               func(w http.ResponseWriter, r *http.Request) http_render.Render
	RenderFailed             func(w http.ResponseWriter, r *http.Request, render *http_render.Render, err error)
}

func (this *RequestRenderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ext := path.Ext(r.URL.Path); ext != "" {
		return
	}
	var Render http_render.Render
	if this.RenderFunc != nil {
		Render = this.RenderFunc(w, r)
	} else {
		Render = this.GetOrCreateRender().Option(ropt.Request(r))
	}
	rd := Render.Option(ropt.FileNames(r.URL.Path))
	rd.Status = http.StatusOK
	if this.RenderOrNotFoundDisabled {
		if err := rd.Render(w); err != nil && !os.IsNotExist(err) && !strings.Contains(err.Error(), "not a directory") {
			if this.RenderFailed == nil {
				panic(err)
			}
			this.RenderFailed(w, r, &rd, err)
		}
	} else {
		rd.MustRenderOrNotFound(w)
	}
}

func (this *RequestRenderHandler) GetOrCreateRender() http_render.Render {
	if this.RootDir == "" {
		this.RootDir = "www"
	}
	return this.Render.Option(ropt.Dir(this.RootDir))
}
