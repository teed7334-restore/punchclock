package models

import (
	"log"

	db "../database"
)

//NoNeedCheckinList 免打卡人員列表
type NoNeedCheckinList struct {
	ID       int `gorm:"AUTO_INCREMENT"`
	MemberID string
}

//GetNoNeedCheckinList 取得免打卡人員列表
func GetNoNeedCheckinList() []*NoNeedCheckinList {
	list := []*NoNeedCheckinList{}
	err := db.Db.Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}
