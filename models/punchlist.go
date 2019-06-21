package models

import (
	"log"
	"strings"
	"time"

	db "github.com/teed7334-restore/punchclock/database"
	"github.com/teed7334-restore/punchclock/env"
)

//PunchList 卡鐘細項資料表結構
type PunchList struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	PunchTime time.Time
	DoorNo    string
	CardNo    string
	Identify  string
}

//DailyPunchList 當天需打卡員工記錄
type DailyPunchList struct {
	OnWorkTime  time.Time
	OffWorkTime time.Time
	Identify    string
	Firstname   string
	Lastname    string
	Email       string
}

//AddPunchList 新增卡鐘細項
func AddPunchList(p *PunchList) {
	err := db.Db.Create(&p).Error
	if err != nil {
		log.Fatal(err)
	}
}

//GetDailyPunchList 取得當天需打卡員工記錄
func GetDailyPunchList(checkTime time.Time) []*DailyPunchList {
	cfg := env.GetEnv()
	now := checkTime.Format(cfg.TimeFormat)
	nowArr := strings.Split(now, " ")
	begin := nowArr[0] + " 00:00:00"
	end := nowArr[0] + " 23:59:59"
	sql := `
		SELECT 
			min(pl.punch_time) AS on_work_time, 
			max(pl.punch_time) AS off_work_time, 
			pl.identify, 
			u.firstname, 
			u.lastname, 
			u.email
		FROM 
			punch_lists pl
		INNER JOIN 
			users u ON u.identifier = pl.identify AND 
			pl.punch_time >= ? AND 
			pl.punch_time <= ? AND
			pl.identify NOT IN (SELECT member_id FROM no_need_checkin_lists)
		GROUP BY 
			pl.identify, 
			u.firstname, 
			u.lastname, 
			u.email, 
			u.manager
	`
	list := []*DailyPunchList{}
	db.Db.Raw(sql, begin, end).Scan(&list)
	return list
}
