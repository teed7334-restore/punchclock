package models

import (
	"log"
	"time"

	db "punchclock/database"
)

//Attendance 出勤記錄表
type Attendance struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	Identify  string
	Late      bool
	Early     bool
	Unchecked bool
	CreateAt  time.Time
}

//AddAttendance 新增出勤記錄
func AddAttendance(a *Attendance) {
	err := db.Db.Create(&a).Error
	if err != nil {
		log.Fatal(err)
	}
}
