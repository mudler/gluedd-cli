package types

import (
	"github.com/mudler/gluedd-cli/pkg/assetstore"
	"github.com/mudler/gluedd/pkg/resource"
	"encoding/hex"
	"math/rand"
	"fmt"
	"path/filepath"
)

func NewJpegStreamer(url string, baseurl string, local string) resource.Resource {
	return &JpegStreamer{StreamUrl: url, BaseUrl: baseurl, Store: local}
}

func TempFileName(prefix, suffix string) string {
    randBytes := make([]byte, 16)
    rand.Read(randBytes)
    return prefix+hex.EncodeToString(randBytes)+suffix
}

type JpegStreamer struct {
	StreamUrl string
	BaseUrl string
	Store string
}

func (l *JpegStreamer) Listen() chan string {
	as:=assetstore.NewAssetStore(l.Store)
	go as.Run()
	files := make(chan string)
	go func() {
		for {
asset:=TempFileName("predict",".jpg")
			err:=assetstore.DownloadFile(filepath.Join(l.Store,asset), l.StreamUrl)
			if err != nil {
				fmt.Println("error downloading stream file")
			}

			files <- l.BaseUrl+asset

		}
	}()

	return files

}
