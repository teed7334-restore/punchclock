package models

import (
	"log"
	"time"

	db "../database"
)

//Leave 請假記錄
type Leave struct {
	ID            int `gorm:"AUTO_INCREMENT"`
	Startdate     time.Time
	Enddate       time.Time
	Status        int
	Employee      int
	Cause         string
	Startdatetype string
	Enddatetype   string
	Duration      float32
	Type          int
	Comments      string
}

//GetLeaveMemberList 取得請假人員清單
func GetLeaveMemberList(identifies []string, checkTime string) []*Leave {
	list := []*Leave{}
	err := db.Db.Joins("INNER JOIN users ON users.id = leaves.employee").Where("leaves.startdate <= ? AND leaves.enddate >= ?", checkTime, checkTime).Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}
