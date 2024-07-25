package logs

import (
	"os"
	"time"
)

func _writeLog(category string, msg string) error {
	file, err := os.OpenFile("./conf/"+category+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	size := fileInfo.Size()
	if size > 50*1024*1024*1024 {
		file.Truncate(25 * 1024 * 1024 * 1024)
	}

	// 将写入置入到顶部
	file.WriteString(getNow() + "    " + msg + " \n\n")
	//file.WriteAt([]byte(getNow()+"    "+msg+" \n\n"), 0)
	//及时关闭file句柄
	defer file.Close()
	////写入文件时，使用带缓存的 *Writer
	//write := bufio.NewWriter(file)
	//write.WriteString(getNow() + "    " + msg + " \n\n")
	////Flush将缓存的文件真正写入到文件中
	//write.Flush()
	return nil
}

func WriteLog(category string, msg string) error {
	err := _writeLog(category, msg)
	return err
}

func getNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
