package ropt

import (
	"context"
	"net/http"

	"github.com/moisespsena-go/http-render"
)

type (
	r                    = http_render.Render
	assetFinder          = http_render.AssetFactory
	Option               = http_render.Option
	directoryAssetFinder = http_render.DirectoryAssetFactory
)

func AssetFinder(finder assetFinder) Option {
	return Option(func(r r) r {
		r.AssetFactory = finder
		return r
	})
}

func FuncMap(funcMap ...map[string]interface{}) Option {
	return Option(func(r r) r {
		r.TemplateFuncMap = r.TemplateFuncMap.Merge(funcMap...)
		return r
	})
}

func Data(data interface{}) Option {
	return Option(func(r r) r {
		r.Data = data
		return r
	})
}

func FileNames(fileNames ...string) Option {
	return Option(func(r r) r {
		r.FileNames = fileNames
		return r
	})
}

func AppendFileNames(fileNames ...string) Option {
	return Option(func(r r) r {
		r.FileNames = append(r.FileNames, fileNames...)
		return r
	})
}

func DirectoryIndexDisabled() Option {
	return Option(func(r r) r {
		r.DirectoryIndexDisabled = true
		return r
	})
}

func DirectoryIndexEnabled() Option {
	return Option(func(r r) r {
		r.DirectoryIndexDisabled = false
		return r
	})
}

func Status(value int) Option {
	return Option(func(r r) r {
		r.Status = value
		return r
	})
}

func Dir(pth string) Option {
	return AssetFinder(&directoryAssetFinder{pth})
}

func Request(req *http.Request) Option {
	return func(r http_render.Render) r {
		r.Request = req
		return r
	}
}

func Context(ctx context.Context) Option {
	return func(r http_render.Render) r {
		r.Request = r.Request.WithContext(ctx)
		return r
	}
}

func ContextValue(key, value interface{}) Option {
	return func(r http_render.Render) r {
		r.Request = r.Request.WithContext(context.WithValue(r.Request.Context(), key, value))
		return r
	}
}
