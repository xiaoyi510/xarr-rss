package fileUtil

import (
	"XArr-Rss/util/logsys"
	"fmt"
	"io"
	"os"
)

func GetFilePathFromList(path string, imgFiles []string) string {
	var posterFileName string
	// 寻找文件是否存在
	for _, imgName := range imgFiles {
		stat, err := os.Stat(path + imgName)
		if err != nil {
			continue
		}
		posterFileName = stat.Name()
		break
	}
	logsys.Info("找到图片:%v", "获取缩略图", posterFileName)
	// 判断文件是否已找到
	return posterFileName
}

// CopyFile copies a file from src to dst. If src and dst files exist,and are
// the same,then return success. Otherise,attempt to create a hard link
// between the two files. If that fail,copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cAnnot copy non-regular files (e.g.,directories,// symlinks,devices,etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular desTination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// desTination file exists,all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
