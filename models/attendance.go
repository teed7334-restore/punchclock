package models

import (
	"log"
	"os"
	"strings"
	"time"

	db "github.com/teed7334-restore/punchclock/database"
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

//AttendanceRecord 出勤記錄表
type AttendanceRecord struct {
	Identify    string
	Firstname   string
	Lastname    string
	Email       string
	OnWorkTime  time.Time
	OffWorkTime time.Time
	Late        bool
	Early       bool
	Manager     int
}

//AddAttendance 新增出勤記錄
func AddAttendance(a *Attendance) {
	err := db.Db.Create(&a).Error
	if err != nil {
		log.Println(err)
	}
}

//GetNoCheckinMember 取得沒打卡員工列表
func GetNoCheckinMember(now time.Time) []*AttendanceRecord {
	list := []*AttendanceRecord{}
	timeFormat := os.Getenv("timeFormat")
	nowStr := now.Format(timeFormat)
	nowArr := strings.Split(nowStr, " ")
	sql := `
		SELECT 
			u.identifier AS identify, 
			u.firstname, 
			u.lastname, 
			u.email, 
			min(pl.punch_time) AS on_work_time, 
			max(pl.punch_time) AS off_work_time, 
			a.late, 
			a.early,
			u.manager
		FROM 
			attendances a
		INNER JOIN 
			users u ON a.identify = u.identifier AND 
			a.create_at >= ? AND 
			a.create_at <= ? AND 
			a.unchecked = 1
		INNER JOIN 
			punch_lists pl ON pl.identify = u.identifier AND 
			pl.punch_time >= ? AND 
			pl.punch_time <= ?
		GROUP BY 
			u.identifier, 
			u.firstname, 
			u.lastname, 
			u.email, 
			a.late, 
			a.early, 
			u.manager`
	begin := nowArr[0] + " 00:00:00"
	end := nowArr[0] + " 23:59:59"
	err := db.Db.Raw(sql, begin, end, begin, end).Scan(&list).Error
	if err != nil {
		log.Println(err)
	}
	return list
}

//GetAttendance 取得出勤記錄表
func GetAttendance(now time.Time) []*AttendanceRecord {
	list := []*AttendanceRecord{}
	timeFormat := os.Getenv("timeFormat")
	nowStr := now.Format(timeFormat)
	nowArr := strings.Split(nowStr, " ")
	sql := `
		SELECT 
			u.identifier AS identify, 
			u.firstname, 
			u.lastname, 
			u.email, 
			min(pl.punch_time) AS on_work_time, 
			max(pl.punch_time) AS off_work_time, 
			a.late, 
			a.early,
			u.manager
		FROM 
			attendances a
		INNER JOIN 
			users u ON a.identify = u.identifier AND 
			a.create_at >= ? AND 
			a.create_at <= ? AND 
			(a.late = 1 OR a.early = 1)
		INNER JOIN 
			punch_lists pl ON pl.identify = u.identifier AND 
			pl.punch_time >= ? AND 
			pl.punch_time <= ?
		GROUP BY 
			u.identifier, 
			u.firstname, 
			u.lastname, 
			u.email, 
			a.late, 
			a.early, 
			u.manager`
	begin := nowArr[0] + " 00:00:00"
	end := nowArr[0] + " 23:59:59"
	err := db.Db.Raw(sql, begin, end, begin, end).Scan(&list).Error
	if err != nil {
		log.Println(err)
	}
	return list
}
