package models

import (
	"log"

	db "github.com/teed7334-restore/punchclock/database"
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
		log.Println(err)
	}
	return list
}
