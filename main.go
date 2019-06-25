package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/teed7334-restore/homekeeper/beans"
	homekeeper "github.com/teed7334-restore/homekeeper/controllers"
	"github.com/teed7334-restore/punchclock/base"
	db "github.com/teed7334-restore/punchclock/database"
	"github.com/teed7334-restore/punchclock/env"
	hook "github.com/teed7334-restore/punchclock/hooks"
	model "github.com/teed7334-restore/punchclock/models"
)

var cfg = env.GetEnv()
var mailLists map[int]string

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

//processDutyData 處理上班資料
func processDutyData(createAt time.Time) {
	cfg := env.GetEnv()
	checkinList := model.GetDailyPunchList(createAt)
	for _, item := range checkinList {
		early := false
		late := false
		onWorkTimeStr := item.OnWorkTime.Format(cfg.TimeFormat)
		onWorkTimeArr := strings.Split(onWorkTimeStr, " ")
		diffDay, diffHour, _ := homekeeper.CallCalcTime(item.OnWorkTime, item.OffWorkTime)
		checkTimeStr := onWorkTimeArr[0] + " " + cfg.CheckTime + ":00"
		checkTime, _ := time.ParseInLocation(cfg.TimeFormat, checkTimeStr, time.Local)
		if checkTime.Before(item.OnWorkTime) { //遲到
			late = true
		}
		diffHour = diffDay*8 + diffHour
		if cfg.DailyWorkHours > diffHour { //早退
			early = true
		}
		list := model.Attendance{Identify: item.Identify, Late: late, Early: early, Unchecked: false, CreateAt: createAt}
		model.AddAttendance(&list)
	}
	noCheckinList := model.GetNoCheckInMember(createAt)
	for _, item := range noCheckinList { //當日未打卡
		list := model.Attendance{Identify: item.Identifier, Late: false, Early: false, Unchecked: true, CreateAt: createAt}
		model.AddAttendance(&list)
	}
}

//combinMemberListMapping 組合員工資料表
func combinMemberListMapping() map[int]string {
	result := model.GetMemberList()
	list := make(map[int]string)
	for _, member := range result {
		key := member.ID
		value := member.Email
		list[key] = value
	}
	return list
}

//noNeedCheckinList 不用打卡員工列表
func noNeedCheckinList() []string {
	denyList := model.GetNoNeedCheckinList()
	list := []string{}
	for _, v := range denyList {
		list = append(list, v.MemberID)
	}
	return list
}

//getHoliday 取得假日資料
func getHoliday() map[string]int {
	list := model.GetHoliday()
	isHoliday := make(map[string]int)
	for _, value := range list {
		dateStr := value.Date.Format(cfg.TimeFormat)
		dateArr := strings.Split(dateStr, " ")
		key := dateArr[0]
		isHoliday[key] = value.IsHoliday
	}
	return isHoliday
}

//skipNoNeedCheckinDays 跳過免打卡節日
func skipNoNeedCheckinDays(isHoliday map[string]int, checkTime string) bool {
	_, ok := isHoliday[checkTime]
	if cfg.Filters.SkipNoNeedCheckinDays && ok { //略過免打卡日期
		return true
	}
	return false
}

//checkTodayHavePunchData 檢查今日打卡資料是否齊全
func checkTodayHavePunchData(isHoliday map[string]int) bool {
	now := time.Now().Format(cfg.TimeFormat)
	today := strings.Split(now, " ")[0]
	checkTime := today
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

//sendMailForWorningCheckinMember 針對打卡有問題員工發送通知信
func sendMailForWorningCheckinMember() {
	cfg := env.GetEnv()
	now := time.Now()
	nowStr := now.Format(cfg.TimeFormat)
	nowArr := strings.Split(nowStr, " ")
	list := model.GetAttendance(now)
	mailLists = combinMemberListMapping()
	for _, item := range list {
		name := item.Firstname
		onWorkTime := item.OnWorkTime.Format(cfg.TimeFormat)
		offWorkTime := item.OffWorkTime.Format(cfg.TimeFormat)
		if item.Late {
			email := item.Email
			subject := "遲到通知"
			content := "親愛的 " + name + " 您於今日 " + onWorkTime + " 才到，目前系統還不會將假期也拉過來計算，如果您己有請假，可自動忽略此次訊息"
			cc := ""
			_, ok := mailLists[item.Manager]
			if ok {
				cc = mailLists[item.Manager]
			}
			sendMail := &beans.SendMail{To: email, Cc: cc, Subject: subject, Content: content}
			hook.SendMail(sendMail)
		}
		if item.Early {
			email := item.Email
			subject := "早退通知"
			content := "親愛的 " + name + " 您於今日 " + offWorkTime + " 離開，目前系統還不會將假期也拉過來計算，如果您己有請假，可自動忽略此次訊息"
			cc := ""
			_, ok := mailLists[item.Manager]
			if ok {
				cc = mailLists[item.Manager]
			}
			sendMail := &beans.SendMail{To: email, Cc: cc, Subject: subject, Content: content}
			hook.SendMail(sendMail)
		}
	}
	list = model.GetNoCheckinMember(now)
	for _, item := range list {
		email := item.Email
		name := item.Firstname
		subject := "未打卡通知"
		content := "親愛的 " + name + " :\r\n" + "您於 " + nowArr[0] + " 尚未打卡"
		cc := ""
		_, ok := mailLists[item.Manager]
		if ok {
			cc = mailLists[item.Manager]
		}
		sendMail := &beans.SendMail{To: email, Cc: cc, Subject: subject, Content: content}
		hook.SendMail(sendMail)
	}
}

//processPunchData 處理卡鐘資料
func processPunchData() int {
	files := base.GetFileList()
	isHoliday := getHoliday() //取得不用打卡日期
	havePunchData := checkTodayHavePunchData(isHoliday)
	deny := []string{}              //取得不用打卡員工列表
	denyList := map[string]string{} //不用打卡員工列表List形式
	if "production" == cfg.Env {    //當使用的HRM系統為Jorani時，才標記未打卡員工
		deny = noNeedCheckinList()
		denyList = map[string]string{}
		for _, v := range deny {
			denyList[v] = "1"
		}
	}
	if !havePunchData {
		subject := "無有效卡鐘檔通知"
		content := "今日沒有可用的卡鐘檔"
		for _, mail := range cfg.AlertMail {
			sendMail := &beans.SendMail{To: mail, Subject: subject, Content: content}
			hook.SendMail(sendMail)
		}
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
			searchTime = searchTime + " 00:00:00"
			createAt, _ := time.ParseInLocation(cfg.TimeFormat, searchTime, time.Local)
			processDutyData(createAt)
		}
	}
	return 1
}

func main() {
	mailLists = combinMemberListMapping()
	if 3 != hook.UpdateHoliday() {
		log.Fatal("Hook Error")
	}
	if 1 != processPunchData() {
		log.Fatal("Punch Data Error")
	}
	if "production" == cfg.Env { //當使用的HRM系統為Jorani時，才標記未打卡員工
		sendMailForWorningCheckinMember()
	}
	defer db.Db.Close()
	fmt.Println(1)
}
