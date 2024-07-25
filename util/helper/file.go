package helper

import (
	"path"
)

type FileHelper struct {
}

func (this FileHelper) GetFileExtension(filename string) string {
	return path.Ext(filename)
}
