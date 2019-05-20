package models

import (
	"log"

	db "punchclock/database"
)

//PunchLog 卡鐘檔記錄資料表結構
type PunchLog struct {
	ID   int `gorm:"AUTO_INCREMENT"`
	Name string
}

//AddPunchLog 寫入卡鐘檔記錄
func AddPunchLog(fileName string) {
	err := db.Db.Create(&PunchLog{Name: fileName}).Error
	if err != nil {
		log.Fatal(err)
	}
}

//CheckPunchLog 檢查卡鐘檔記錄
func CheckPunchLog(fileName string) []*PunchLog {
	list := []*PunchLog{}
	err := db.Db.Where("name = ?", fileName).Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}
