# punchclock
701卡鐘資料匯入程式

## 資料夾結構
base 底層共用Library

beans 用來裝Call API後之ResultObject

data 卡鐘檔

database 資料庫Driver檔

env 系統設定

hooks 呼叫第三方API

models 資料庫存取邏輯

main.go 主程式

## 必須套件
本程式透過Google Protobuf 3產生所需之ResultObject，然Proto 3之後官方不支持Custom Tags，所以還需要多安裝一個寫入retags的套件

```
git clone https://github.com/qianlnk/protobuf.git $GOPATH/src/github.com/golang/protobuf

go install $GOPATH/src/github.com/golang/protobuf/protoc-gen-go
```

## 程式運行原理
本程式會在你的資料庫新增三張表，分別是punch_logs(卡鐘檔記錄資料表), punch_list(卡鐘細項資料表), attendances(出勤記錄表)

程式運行時，會將卡鐘檔的檔名先存到卡鐘檔記錄資料表，防止同樣檔名的資料重覆匯入

之後會讀取卡鐘檔案的內容，每筆照實存到卡鐘細項資料表中

最後透過存入的卡鐘細項資料，生成遲到、早退記錄

本系統透過佇列管理員與服務管理員進行發信之動作，如果沒有要使用發信功能者，可以將發信功能自己注解掉

佇列管理員
https://github.com/teed7334-restore/counter

服務管理員
https://github.com/teed7334-restore/homekeeper

## 程式操作流程
1. 將701 Server匯出的卡鐘檔扔到./data底下，您也可以透過修改原始碼變更路徑
2. 卡鐘檔需以下格式-[日期][時間][門號][卡號][員工編號]，且透過空白做分隔，樣本可參照./data/20190312.txt
3. 將./env/env.swp檔名改成env.go
4. 修改./env/env.go換成你連入的MySQL資料庫與帳號密碼，如果您不是使用[Jorani](https://jorani.org)這套HRM，請將cfg.Env = "production"改成其他任一值
5. 到./beans底下，運行protoc --go_out=plugins=grpc+retag:. *.proto
6. 運行main.go

## 其它
本程式目前有計算曠職功能，並寄發通知信，但需要搭配[Jorani](https://jorani.org)這套HRM系統一起使用

本程式會自動去抓人事行政局休假日資料，寫到holidays資料表，可供公司內部HRM與ERP做參照