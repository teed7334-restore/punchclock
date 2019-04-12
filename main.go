package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	db "./database"
	env "./env"
	hook "./hooks"
	model "./models"
)

var cfg = env.GetEnv()

//getFileList 開啟資料夾中檔案列表
func getFileList() []os.FileInfo {
	files, err := ioutil.ReadDir(cfg.Path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

//getRowData 取得檔案內容
func getRowData(fileName string) *bufio.Scanner {
	txt, err := os.Open(cfg.Path + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(txt)
	return scanner
}

func init() {
	//如果您要修改本程式碼，可以關閉以下方註解以方便測試
	//db.Db.DropTable(&model.PunchLog{}, &model.PunchList{}, &model.Attendance{})
	db.Db.DropTable(&model.Holiday{})
	db.Db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Holiday{}, &model.PunchLog{}, &model.PunchList{}, &model.Attendance{})
}

//combinDutyData 組合上班資料
func combinDutyData(searchTime string) map[string]string {
	daily := model.GetDailyPunchList(searchTime)
	duty := make(map[string]string)
	for _, v := range daily {
		_, ok := duty[v.Identify]
		if ok {
			duty[v.Identify] = duty[v.Identify] + v.PunchTime.Format(cfg.TimeFormat) + ","
		} else {
			duty[v.Identify] = v.PunchTime.Format(cfg.TimeFormat) + ","
		}
	}
	return duty
}

//processDutyData 處理上班資料
func processDutyData(duty map[string]string, createAt time.Time) {
	for k, v := range duty {
		item := strings.Split(v, ",")
		goWork, _ := time.ParseInLocation(cfg.TimeFormat, item[0], time.Local)
		outWork, _ := time.ParseInLocation(cfg.TimeFormat, item[len(item)-2], time.Local)
		item = strings.Split(item[0], " ")
		am, _ := time.ParseInLocation(cfg.TimeFormat, item[0]+" "+cfg.WorkAt+":00", time.Local)
		workTime := outWork.Sub(goWork).Hours() - cfg.LunchTimeHours
		late := false
		early := false
		if cfg.DailyWorkHours >= workTime {
			early = true
		}
		if goWork.After(am) {
			late = true
		}
		list := model.Attendance{Identify: k, Late: late, Early: early, Unchecked: false, CreateAt: createAt}
		model.AddAttendance(&list)
	}
}

//markUnClockMember 標記沒打卡員工
func markUnClockMember(ids map[string]int, searchTime string, createAt time.Time) {
	identifies := []string{}
	for k := range ids {
		identifies = append(identifies, k)
	}
	cmlResult := model.GetClockMemberList(identifies)
	memberIds := []string{}
	uncheckList := make(map[int]string)
	for k := range cmlResult {
		memberIds = append(memberIds, cmlResult[k].Identifier)
		uncheckList[cmlResult[k].ID] = cmlResult[k].Identifier
	}
	lmlResult := model.GetLeaveMemberList(memberIds, searchTime)
	for k := range lmlResult {
		key := lmlResult[k].Employee
		_, ok := uncheckList[key]
		if ok {
			uncheckList[key] = ""
		}
	}
	for k := range uncheckList {
		if "" != uncheckList[k] {
			list := model.Attendance{Identify: uncheckList[k], Late: false, Early: false, Unchecked: true, CreateAt: createAt}
			model.AddAttendance(&list)
		}
	}
}

//processPunchData 處理卡鐘資料
func processPunchData() int {
	files := getFileList()
	for _, f := range files {
		fileName := f.Name()
		list := model.CheckPunchLog(fileName)
		if 0 == len(list) {
			model.AddPunchLog(fileName)
			scanner := getRowData(fileName)
			searchTime := ""
			ids := make(map[string]int)
			for scanner.Scan() { //將文字檔資料寫入資料表
				item := strings.Split(scanner.Text(), " ")
				y := "20" + item[0][0:2]
				m := item[0][2:4]
				d := item[0][4:6]
				h := item[1][0:2]
				i := item[1][2:4]
				s := "00"
				searchTime = fmt.Sprintf("%s-%s-%s", y, m, d)
				checkTime := fmt.Sprintf("%s-%s-%s %s:%s:%s", y, m, d, h, i, s)
				punchTime, _ := time.ParseInLocation(cfg.TimeFormat, checkTime, time.Local)
				list := model.PunchList{PunchTime: punchTime, DoorNo: item[2], CardNo: item[3], Identify: item[4]}
				ids[item[4]] = 1
				model.AddPunchList(&list)
			}
			duty := combinDutyData(searchTime)
			searchTime = searchTime + " 00:00:00"
			createAt, _ := time.ParseInLocation(cfg.TimeFormat, searchTime, time.Local)
			processDutyData(duty, createAt)
			if "production" == cfg.Env { //當使用的HRM系統為Jorani時，才標記未打卡員工
				markUnClockMember(ids, searchTime, createAt)
			}
			defer db.Db.Close()
		}
	}
	return 1
}

func main() {
	channel := make(chan int)
	go func() { channel <- hook.UpdateHoliday() }()
	go func() { channel <- processPunchData() }()
	result := <-channel + <-channel
	fmt.Println(result)
}
