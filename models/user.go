package models

import (
	"log"
	"time"

	db "punchclock/database"
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

//GetClockMemberList 取得沒打卡人員列表
func GetClockMemberList(identifies []string) []*User {
	list := []*User{}
	err := db.Db.Where("identifier NOT IN (?)", identifies).Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}
