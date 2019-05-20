package hooks

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/teed7334-restore/punchclock/base"
	bean "github.com/teed7334-restore/punchclock/beans"
	"github.com/teed7334-restore/punchclock/env"
	model "github.com/teed7334-restore/punchclock/models"
)

var cfg = env.GetEnv()

//UpdateHoliday 連上人事行政局處理回傳之休假資料寫入資料表
func UpdateHoliday() int {
	url := cfg.API.Holiday
	body := base.GetURL(url)
	resultObject := bean.Holiday{}
	json.Unmarshal(body, &resultObject)
	success := resultObject.GetSuccess()
	if success {
		items := resultObject.GetResult().GetRecord()
		for _, value := range items {
			dateSrc := value.GetDate()
			dateArray := strings.Split(dateSrc, "/")
			year := dateArray[0]
			month := dateArray[1]
			if 2 > len(month) {
				month = "0" + month
			}
			day := dateArray[2]
			if 2 > len(day) {
				day = "0" + day
			}
			dateSrc = year + "-" + month + "-" + day + " 00:00:00"
			date, _ := time.ParseInLocation(cfg.TimeFormat, dateSrc, time.Local)
			name := value.GetName()
			holiday := value.GetIsHoliday()
			isHoliday := 0
			if "是" == holiday {
				isHoliday = 1
			}
			holidayCategory := value.GetHolidayCategory()
			description := value.GetDescription()
			list := model.Holiday{Date: date, Name: name, IsHoliday: isHoliday, HolidayCategory: holidayCategory, Description: description}
			model.AddHoliday(&list)
		}
	}
	return 3
}
