package graphite

import (
	"io/ioutil"
	"strings"
)

type PlainTextFile struct {
	path    string
	Name    string
	Ext     string
	Content []byte
}

func StripNameExt(path string) (name, ext string) {
	path_split := strings.Split(path, "/")
	splited := strings.Split(path_split[len(path_split)-1], ".")
	return splited[0], splited[len(splited)-1]
}

func PlainText(path string) (PlainTextFile, error) {
	file := PlainTextFile{path: path}
	file.Name, file.Ext = StripNameExt(file.path)
	content, err := ioutil.ReadFile(file.path)
	if err != nil {
		return PlainTextFile{}, err
	}
	file.Content = content
	return file, nil
}
