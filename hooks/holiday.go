package hooks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/teed7334-restore/punchclock/base"
	"github.com/teed7334-restore/punchclock/beans"
	"github.com/teed7334-restore/punchclock/models"
)

//UpdateHoliday 連上人事行政局處理回傳之休假資料寫入資料表
func UpdateHoliday() int {
	url := os.Getenv("api.holiday")
	timeFormat := os.Getenv("timeFormat")
	body := base.GetURL(url)
	response := &[]beans.Holiday{}
	err := json.Unmarshal(body, response)
	if err == nil {
		for _, item := range *response {
			arr := strings.Split(item.Date, "/")
			year := arr[0]
			month := arr[1]
			if len(month) < 2 {
				month = "0" + month
			}
			day := arr[2]
			if len(day) < 2 {
				day = "0" + day
			}
			date := fmt.Sprintf("%s-%s-%s 00:00:00", year, month, day)
			checkTime, _ := time.ParseInLocation(timeFormat, date, time.Local)
			list := models.Holiday{Date: checkTime, Name: item.Name, IsHoliday: 1, HolidayCategory: "", Description: item.Description}
			models.AddHoliday(&list)
		}
		//將非企業但人事行政局有設為放假的假別取消
		models.SetHolidayForNonBusiness()
	}
	return 3
}
