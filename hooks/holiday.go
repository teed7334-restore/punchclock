package hooks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	env "../env"
	model "../models"
)

var cfg = env.GetEnv()

//getURL 透過HTTP GET取得網頁資料
func getURL(url string) []byte {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	result, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(result.Body)
	defer result.Body.Close()
	return body
}

//jsonDecode 將網頁回傳之文字資料轉成JSON格式
func jsonDecode(body []byte) interface{} {
	var resultObject interface{}
	json.Unmarshal(body, &resultObject)
	return resultObject
}

//UpdateHoliday 連上人事行政局處理回傳之休假資料寫入資料表
func UpdateHoliday() {
	url := "http://data.ntpc.gov.tw/api/v1/rest/datastore/382000000A-000077-002"
	body := getURL(url)
	resultObject := jsonDecode(body)
	success := resultObject.(map[string]interface{})["success"].(bool)
	if success {
		items := resultObject.(map[string]interface{})["result"].(map[string]interface{})["records"].([]interface{})
		for _, value := range items {
			dateSrc := value.(map[string]interface{})["date"].(string)
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
			name := value.(map[string]interface{})["name"].(string)
			holiday := value.(map[string]interface{})["isHoliday"].(string)
			isHoliday := false
			if "是" == holiday {
				isHoliday = true
			}
			holidayCategory := value.(map[string]interface{})["holidayCategory"].(string)
			description := value.(map[string]interface{})["description"].(string)
			list := model.Holiday{Date: date, Name: name, IsHoliday: isHoliday, HolidayCategory: holidayCategory, Description: description}
			model.AddHoliday(&list)
		}
	}
}
