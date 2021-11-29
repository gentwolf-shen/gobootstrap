package embed

import (
	"io/fs"
	"io/ioutil"
	"os"
)

type LocalEmbedFile struct {
}

func (f LocalEmbedFile) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (f LocalEmbedFile) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (f LocalEmbedFile) ReadFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}
