package hooks

import (
	"encoding/json"
	"os"

	"github.com/teed7334-restore/homekeeper/beans"
	"github.com/teed7334-restore/punchclock/base"
)

//SendMail 串接佇列管理員進行發信
func SendMail(sendMail *beans.SendMail) {
	counter := os.Getenv("api.counter")
	params, _ := json.Marshal(sendMail)
	url := counter + "/Mail/SendMail"
	base.PostURL(url, []byte(params))
}
