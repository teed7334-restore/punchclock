package models

import (
	"log"
	"time"

	db "github.com/teed7334-restore/punchclock/database"
)

//Holiday 假期資料表
type Holiday struct {
	ID              int `gorm:"AUTO_INCREMENT"`
	Date            time.Time
	Name            string
	IsHoliday       int
	HolidayCategory string
	Description     string
}

//AddHoliday 新增假期資料
func AddHoliday(h *Holiday) {
	err := db.Db.Create(&h).Error
	if err != nil {
		log.Println(err)
	}
}

//GetHoliday 取得休假日資料
func GetHoliday() []*Holiday {
	list := []*Holiday{}
	err := db.Db.Where("is_holiday = ?", 1).Order("date desc").Find(&list).Error
	if err != nil {
		log.Println(err)
	}
	return list
}

//SetHolidayForNonBusiness 將非企業但人事行政局有設為放假的假別取消
func SetHolidayForNonBusiness() {
	sql := `
		UPDATE 
			holidays 
		SET 
			is_holiday = 0
		WHERE 
			name IN ('軍人節', '教師節') AND 
			holiday_category IN ('特定節日')
	`
	db.Db.Exec(sql)
}
