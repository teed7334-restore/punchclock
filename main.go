package main

import (
	"fmt"
	"strings"
	"time"

	"punchclock/base"

	db "punchclock/database"

	"punchclock/env"

	hook "punchclock/hooks"

	model "punchclock/models"
)

var cfg = env.GetEnv()

func init() {
	if db.Db.HasTable(&model.PunchLog{}) && db.Db.HasTable(&model.PunchList{}) && db.Db.HasTable(&model.Attendance{}) && db.Db.HasTable(&model.NoNeedCheckinList{}) {
		//如果您要修改本程式碼，可以關閉以下方註解以方便測試
		//db.Db.DropTable(&model.PunchLog{}, &model.PunchList{}, &model.Attendance{})
		//db.Db.DropTable(&model.NoNeedCheckinList{})
	}
	if db.Db.HasTable(&model.Holiday{}) {
		db.Db.DropTable(&model.Holiday{})
	}
	db.Db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Holiday{}, &model.PunchLog{}, &model.PunchList{}, &model.Attendance{}, &model.NoNeedCheckinList{})
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
func markUnClockMember(ids map[string]int, searchTime string, createAt time.Time) map[string]string {
	identifies := []string{}
	for k := range ids {
		identifies = append(identifies, k)
	}
	cmlResult := model.GetClockMemberList(identifies) //取得需打卡員工列表
	memberIds := []string{}
	uncheckList := make(map[int]string)
	sendCheckList := make(map[string]string)
	deny := noNeedCheckinList() //取得不用打卡員工列表
	for k := range cmlResult {
		id := cmlResult[k].ID
		email := cmlResult[k].Email
		lastName := cmlResult[k].Lastname
		firstName := cmlResult[k].Firstname
		memberID := cmlResult[k].Identifier
		dueDate := cmlResult[k].Datehired.Format(cfg.TimeFormat)
		if skipNoNeedCheckinMembers(deny, memberID) {
			continue
		}
		if skipUnsuccessful(searchTime, dueDate) {
			fmt.Println(lastName)
			continue
		}
		memberIds = append(memberIds, memberID)
		uncheckList[id] = memberID
		sendCheckList[email] = lastName + " " + firstName
	}
	lmlResult := model.GetLeaveMemberList(memberIds, searchTime) //取得員工請假列表
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
	return sendCheckList
}

//noNeedCheckinList 不用打卡員工列表
func noNeedCheckinList() map[string]int {
	denyList := model.GetNoNeedCheckinList()
	list := make(map[string]int)
	for k := range denyList {
		key := denyList[k].MemberID
		list[key] = 1
	}
	return list
}

//getHoliday 取得假日資料
func getHoliday() map[string]int {
	list := model.GetHoliday()
	isHoliday := make(map[string]int)
	for _, value := range list {
		key := value.Date.Format(cfg.TimeFormat)
		isHoliday[key] = value.IsHoliday
	}
	return isHoliday
}

//skipNoNeedCheckinMembers 跳過免打卡人員
func skipNoNeedCheckinMembers(deny map[string]int, identify string) bool {
	_, ok := deny[identify]
	if cfg.Filters.SkipNoNeedCheckinMembers && ok { //略過不用打卡名單
		return true
	}
	return false
}

//skipNoNeedCheckinDays 跳過免打卡節日
func skipNoNeedCheckinDays(isHoliday map[string]int, checkTime string) bool {
	_, ok := isHoliday[checkTime]
	if cfg.Filters.SkipNoNeedCheckinDays && ok { //略過免打卡日期
		return true
	}
	return false
}

//skipUnsuccessful 跳過未到職員工
func skipUnsuccessful(checkTime string, dueDate string) bool {
	check, _ := time.ParseInLocation(cfg.TimeFormat, checkTime, time.Local)
	due, _ := time.ParseInLocation(cfg.TimeFormat, dueDate, time.Local)
	if check.After(due) {
		return true
	}
	return false
}

//checkTodayHavePunchData 檢查今日打卡資料是否齊全
func checkTodayHavePunchData(isHoliday map[string]int) bool {
	now := time.Now().Format(cfg.TimeFormat)
	today := strings.Split(now, " ")[0]
	checkTime := today + " 00:00:00"
	dateArr := strings.Split(today, "-")
	y := dateArr[0]
	m := dateArr[1]
	d := dateArr[2]
	checkFileName := y + m + d + ".txt"
	haveFile := base.CheckFile(checkFileName)
	noNeedCheckInDay := skipNoNeedCheckinDays(isHoliday, checkTime)
	if !noNeedCheckInDay && !haveFile {
		return false
	}
	return true
}

//processPunchData 處理卡鐘資料
func processPunchData() int {
	files := base.GetFileList()
	isHoliday := getHoliday() //取得不用打卡日期
	havePunchData := checkTodayHavePunchData(isHoliday)
	if !havePunchData {
		base.SendMail(cfg.AlertMail, "無有效卡鐘檔通知", "今日沒有可用的卡鐘檔")
	}
	for _, f := range files {
		fileName := f.Name()
		list := model.CheckPunchLog(fileName)
		if 0 == len(list) {
			model.AddPunchLog(fileName)
			scanner := base.GetRowData(fileName)
			searchTime := ""
			ids := make(map[string]int)
			for scanner.Scan() { //將文字檔資料寫入資料表
				item := strings.Split(scanner.Text(), " ")
				doorNo := item[2]
				cardNo := item[3]
				identify := item[4]
				y := "20" + item[0][0:2]
				m := item[0][2:4]
				d := item[0][4:6]
				h := item[1][0:2]
				i := item[1][2:4]
				s := "00"
				searchTime = fmt.Sprintf("%s-%s-%s", y, m, d)
				checkTime := fmt.Sprintf("%s-%s-%s %s:%s:%s", y, m, d, h, i, s)
				punchTime, _ := time.ParseInLocation(cfg.TimeFormat, checkTime, time.Local)
				list := model.PunchList{PunchTime: punchTime, DoorNo: doorNo, CardNo: cardNo, Identify: identify}
				ids[identify] = 1
				model.AddPunchList(&list)
			}
			duty := combinDutyData(searchTime)
			onCheckTime := searchTime
			searchTime = searchTime + " 00:00:00"
			createAt, _ := time.ParseInLocation(cfg.TimeFormat, searchTime, time.Local)
			processDutyData(duty, createAt)
			if "production" == cfg.Env && !skipNoNeedCheckinDays(isHoliday, searchTime) { //當使用的HRM系統為Jorani時，才標記未打卡員工
				unClockMembers := markUnClockMember(ids, searchTime, createAt)
				for email, name := range unClockMembers {
					to := []string{email}
					message := "親愛的 " + name + " :\r\n" + "您於 " + onCheckTime + " 尚未打卡"
					base.SendMail(to, "未打卡通知", message)
				}
			}
			defer db.Db.Close()
		}
	}
	return 1
}

func main() {
	if 3 != hook.UpdateHoliday() {
		panic("Hook Error")
	}
	if 1 != processPunchData() {
		panic("Punch Data Error")
	}
	fmt.Println(1)
}
