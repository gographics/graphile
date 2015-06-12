package graphite

import (
	"io/ioutil"
)

type TextFile struct {
	path    string
	Name    string
	Ext     string
	Content []byte
}

func OpenText(path string) (TextFile, error) {
	file := TextFile{path: path}
	file.Name, file.Ext = PathStrip(file.path)
	content, err := ioutil.ReadFile(file.path)
	if err != nil {
		return TextFile{}, err
	}
	file.Content = content
	return file, nil
}
