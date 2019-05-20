package models

import (
	"log"
	"time"

	db "github.com/teed7334-restore/punchclock/database"
)

//PunchList 卡鐘細項資料表結構
type PunchList struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	PunchTime time.Time
	DoorNo    string
	CardNo    string
	Identify  string
}

//AddPunchList 新增卡鐘細項
func AddPunchList(p *PunchList) {
	err := db.Db.Create(&p).Error
	if err != nil {
		log.Fatal(err)
	}
}

//GetDailyPunchList 取得當天打卡記錄
func GetDailyPunchList(checkTime string) []*PunchList {
	list := []*PunchList{}
	begin := checkTime + " 00:00:00"
	end := checkTime + " 23:59:59"
	err := db.Db.Where("punch_time >= ? AND punch_time <= ?", begin, end).Order("punch_time ASC").Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}
