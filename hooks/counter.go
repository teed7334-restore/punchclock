package hooks

import (
	"encoding/json"

	"github.com/teed7334-restore/punchclock/base"
	"github.com/teed7334-restore/punchclock/beans"
	"github.com/teed7334-restore/punchclock/env"
)

//SendMail 串接佇列管理員進行發信
func SendMail(sendMail *beans.SendMail) {
	cfg := env.GetEnv()
	params, _ := json.Marshal(sendMail)
	url := cfg.API.Counter + "/Mail/SendMail"
	base.PostURL(url, []byte(params))
}
