package models

import (
	"log"
	"time"

	db "../database"
)

//Holiday 假期資料表
type Holiday struct {
	ID              int `gorm:"AUTO_INCREMENT"`
	Date            time.Time
	Name            string
	IsHoliday       bool
	HolidayCategory string
	Description     string
}

//AddHoliday 新增假期資料
func AddHoliday(h *Holiday) {
	err := db.Db.Create(&h).Error
	if err != nil {
		log.Fatal(err)
	}
}
