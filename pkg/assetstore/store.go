// Copyright © 2019 Ettore Di Giacinto <mudler@gentoo.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

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
	m.Run()
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
