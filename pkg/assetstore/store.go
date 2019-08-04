package assetstore

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/sahilm/fuzzy"
	macaron "gopkg.in/macaron.v1"
)

type Store interface {
	Run()
}

type AssetStore struct {
	Directory string
}

func NewAssetStore(dir string) Store {
	return &AssetStore{Directory: dir}
}

func (a *AssetStore) Run() {
	m := macaron.Classic()
	m.Get("/file/:search", func(ctx *macaron.Context) {
		f, err := ServeFuzzy(ctx.Params(":search"), a.Directory)
		if err != nil {
			return
		}
		fi, err := f.Stat()
		if err != nil {
			return
		}
		http.ServeContent(ctx.Resp, ctx.Req.Request, ctx.Req.URL.Path, fi.ModTime(), f)
	})
	go m.Run()
}

func ServeFuzzy(search, dir string) (*os.File, error) {

	f := []string{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, v := range files {
		if !v.IsDir() {

			f = append(f, v.Name())
		}
	}
	matches := fuzzy.Find(search, f)
	if len(matches) == 0 {
		return nil, errors.New("No Matches")
	}

	return os.Open(path.Join(dir, matches[0].Str))

}
