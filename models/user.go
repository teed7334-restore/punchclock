package models

import (
	"log"
	"strings"
	"time"

	db "github.com/teed7334-restore/punchclock/database"
	"github.com/teed7334-restore/punchclock/env"
)

//User 人員資料表
type User struct {
	ID           int `gorm:"AUTO_INCREMENT"`
	Firstname    string
	Lastname     string
	Login        string
	Email        string
	Password     string
	Role         int
	Manager      int
	Country      int
	Organization int
	Contract     int
	Position     int
	Datehired    time.Time
	Identifier   string
	Language     string
	LdapPath     string
	Active       int
	Timezone     string
	Calendar     string
}

//GetNoCheckInMember 取得未打卡員工列表
func GetNoCheckInMember(checkTime time.Time) []*User {
	cfg := env.GetEnv()
	now := checkTime.Format(cfg.TimeFormat)
	nowArr := strings.Split(now, " ")
	begin := nowArr[0] + " 00:00:00"
	end := nowArr[0] + " 23:59:59"
	sql := `
		SELECT 
			*
		FROM 
			users
		WHERE 
			identifier NOT IN (
				SELECT 
					member_id 
				FROM 
					no_need_checkin_lists) AND 
			identifier NOT IN (
				SELECT 
					identify
				FROM 
					punch_lists pl
				WHERE 
					punch_time >= ? AND 
					punch_time <= ?)
	`
	list := []*User{}
	db.Db.Raw(sql, begin, end).Scan(&list)
	return list
}

//GetMemberList 取得全部員工列表
func GetMemberList() []*User {
	list := []*User{}
	err := db.Db.Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}
