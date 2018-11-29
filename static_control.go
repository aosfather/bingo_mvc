package bingo_mvc

import (
	"fmt"
	utils "github.com/aosfather/bingo_utils"
	"mime"
	"os"
	"strings"
)

type staticController struct {
	Controller
	staticDir string
	log       utils.Log
}

func (this *staticController) GetParameType(method string) interface{} {
	return &StaticResource{}

}
func (this *staticController) Get(c Context, p interface{}) (interface{}, BingoError) {
	if resource, ok := p.(*StaticResource); ok {
		this.log.Debug("static resource %s,%s", resource.Type, resource.Uri)
		var view StaticView
		var fileDir string
		fileDir, view.Name, view.Media = parseUri(resource.Uri)

		var filePath string = this.staticDir
		if filePath != "" {
			filePath = filePath + "/"
		}
		if fileDir != "" {
			filePath = filePath + fileDir + "/"
		}
		fileRealPath := filePath + view.Name
		fmt.Print(fileRealPath)

		if utils.IsFileExist(fileRealPath) {
			fi, err := os.Open(fileRealPath)
			if err != nil {
				this.log.Debug(err.Error())
			} else {
				view.Reader = fi
				return view, nil
			}

		}

	}
	return nil, utils.CreateError(Code_NOT_FOUND, "bingo! The uri not found!")

}

func parseUri(uri string) (dir string, name string, media string) {
	fixIndex := strings.LastIndex(uri, ".")
	lastUrlIndex := strings.LastIndex(uri, "/")
	dir = ""
	if lastUrlIndex > 0 {
		dir = string([]byte(uri)[1:lastUrlIndex])
		dir = strings.Replace(dir, "../", "_", -1)
	}

	if lastUrlIndex < 0 {
		lastUrlIndex = 0
	}

	if fixIndex < 0 {
		fixIndex = len(uri)
	}
	var fileSufix string
	querySufixIndex := strings.LastIndex(uri, "?")
	if querySufixIndex > 0 && fixIndex < querySufixIndex {
		fileSufix = string([]byte(uri)[fixIndex:querySufixIndex])
		name = string([]byte(uri)[lastUrlIndex+1 : querySufixIndex])
	} else {
		fileSufix = string([]byte(uri)[fixIndex:])
		name = string([]byte(uri)[lastUrlIndex+1:])
	}
	fmt.Println(fileSufix)
	return dir, name, getMedia(fileSufix)

}

func getMedia(fileFix string) string {
	media := mime.TypeByExtension(fileFix)
	if media == "" {

	}
	return media
}
