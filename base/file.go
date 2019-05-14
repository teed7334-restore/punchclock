package base

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
)

//GetFileList 開啟資料夾中檔案列表
func GetFileList() []os.FileInfo {
	files, err := ioutil.ReadDir(cfg.Path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

//CheckFile 檢查檔案是否存在
func CheckFile(fileName string) bool {
	_, err := os.Stat(cfg.Path + "/" + fileName)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

//GetRowData 取得檔案內容
func GetRowData(fileName string) *bufio.Scanner {
	txt, err := os.Open(cfg.Path + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(txt)
	return scanner
}
