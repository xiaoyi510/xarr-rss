package size

import (
	"fmt"
)

// GetSizeString 获取文件大小
func GetSizeString(size int64) string {
	if size == 0 {
		return ""
	}
	units := []string{" B", " KB", " MB", " GB", " TB", " PB"}
	var i = 0
	size2 := float64(size)
	for i = 0; size2 >= 1024 && i < 5; i++ {
		size2 /= 1024
	}
	format := fmt.Sprintf("%.2f", size2)
	return format + units[i]

}
