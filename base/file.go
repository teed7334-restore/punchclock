package base

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
)

//GetFileList 開啟資料夾中檔案列表
func GetFileList() []os.FileInfo {
	path := os.Getenv("path")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	return files
}

//CheckFile 檢查檔案是否存在
func CheckFile(fileName string) bool {
	path := os.Getenv("path")
	_, err := os.Stat(path + "/" + fileName)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

//GetRowData 取得檔案內容
func GetRowData(fileName string) *bufio.Scanner {
	path := os.Getenv("path")
	txt, err := os.Open(path + "/" + fileName)
	if err != nil {
		log.Println(err)
	}
	scanner := bufio.NewScanner(txt)
	return scanner
}
