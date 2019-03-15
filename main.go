package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

//資料庫使用者
const user = "erp"

//資料庫密碼
const password = "lD9nAKQgYElV2tan"

//資料庫主機位址
const host = "127.0.0.1"

//使用的資料庫
const database = "erp"

//資料庫資料回傳編碼
const charset = "utf8mb4"
const parseTime = "true"
const loc = "Local"

//資料夾路徑
const path = "./data"

//時間格式
const timeFormat = "2006-01-02 15:04:05"

//上班時間
const workAt = "08:30"

//午休時間(小時)
const lunchTimeHours = 1.5

//每天工作時數
const dailyWorkHours = 8

//Database 資料庫連線設定結構
type Database struct {
	User      string
	Password  string
	Host      string
	Database  string
	Charset   string
	ParseTime string
	Loc       string
}

//PunchLog 卡鐘檔記錄資料表結構
type PunchLog struct {
	ID   int `gorm:"AUTO_INCREMENT"`
	Name string
}

//PunchList 卡鐘細項資料表結構
type PunchList struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	PunchTime time.Time
	DoorNo    string
	CardNo    string
	Identify  string
}

//Attendance 出勤記錄表
type Attendance struct {
	ID       int `gorm:"AUTO_INCREMENT"`
	Identify string
	Late     bool
	Early    bool
}

//connectDB 連結資料庫
func connectDB() {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=%s&loc=%s", user, password, host, database, charset, parseTime, loc)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
}

//getFileList 開啟資料夾中檔案列表
func getFileList() []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

//getRowData 取得檔案內容
func getRowData(fileName string) *bufio.Scanner {
	txt, err := os.Open(path + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(txt)
	return scanner
}

//addPunchLog 寫入卡鐘檔記錄
func addPunchLog(fileName string) {
	err = db.Create(&PunchLog{Name: fileName}).Error
	if err != nil {
		log.Fatal(err)
	}
}

//checkPunchLog 檢查卡鐘檔記錄
func checkPunchLog(fileName string) []*PunchLog {
	list := []*PunchLog{}
	err = db.Where("name = ?", fileName).Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}

//addPunchList 新增卡鐘細項
func addPunchList(p *PunchList) {
	err = db.Create(&p).Error
	if err != nil {
		log.Fatal(err)
	}
}

//getDailyPunchList 取得當天打卡記錄
func getDailyPunchList(checkTime string) []*PunchList {
	list := []*PunchList{}
	begin := checkTime + " 00:00:00"
	end := checkTime + " 23:59:59"
	err = db.Where("punch_time >= ? AND punch_time <= ?", begin, end).Order("punch_time ASC").Find(&list).Error
	if err != nil {
		log.Fatal(err)
	}
	return list
}

//addAttendance 新增出勤記錄
func addAttendance(a *Attendance) {
	err := db.Create(&a).Error
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	connectDB()
	//開發時可以取消以下注解
	//db.DropTable(&PunchLog{}, &PunchList{}, &Attendance{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&PunchLog{}, &PunchList{}, &Attendance{})
}

//combinDutyData 組合上班資料
func combinDutyData(searchTime string) map[string]string {
	daily := getDailyPunchList(searchTime)
	duty := make(map[string]string)
	for _, v := range daily {
		_, ok := duty[v.Identify]
		if ok {
			duty[v.Identify] = duty[v.Identify] + v.PunchTime.Format(timeFormat) + ","
		} else {
			duty[v.Identify] = v.PunchTime.Format(timeFormat) + ","
		}
	}
	return duty
}

//processDutyData 處理上班資料
func processDutyData(duty map[string]string) {
	for k, v := range duty {
		item := strings.Split(v, ",")
		goWork, _ := time.ParseInLocation(timeFormat, item[0], time.Local)
		outWork, _ := time.ParseInLocation(timeFormat, item[len(item)-2], time.Local)
		item = strings.Split(item[0], " ")
		am, _ := time.ParseInLocation(timeFormat, item[0]+" "+workAt+":00", time.Local)
		workTime := outWork.Sub(goWork).Hours() - lunchTimeHours
		late := false
		early := false
		if dailyWorkHours >= workTime {
			early = true
		}
		if goWork.After(am) {
			late = true
		}
		list := Attendance{Identify: k, Late: late, Early: early}
		addAttendance(&list)
	}
}

func main() {
	files := getFileList()
	for _, f := range files {
		fileName := f.Name()
		list := checkPunchLog(fileName)
		if 0 == len(list) {
			addPunchLog(fileName)
			scanner := getRowData(fileName)
			searchTime := ""
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
				punchTime, _ := time.ParseInLocation(timeFormat, checkTime, time.Local)
				list := PunchList{PunchTime: punchTime, DoorNo: item[2], CardNo: item[3], Identify: item[4]}
				addPunchList(&list)
			}
			duty := combinDutyData(searchTime)
			processDutyData(duty)
		}
	}
}
