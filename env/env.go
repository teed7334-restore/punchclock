package env

//Config 系統參數
type Config struct {
	Env            string
	Path           string
	TimeFormat     string
	WorkAt         string
	LunchTimeHours float64
	DailyWorkHours float64
	Database       struct {
		Host      string
		User      string
		Password  string
		Database  string
		Charset   string
		ParseTime string
		Loc       string
	}
}

var cfg = &Config{}

func init() {
	cfg.Env = "production"
	cfg.Path = "./data"
	cfg.TimeFormat = "2006-01-02 15:04:05"
	cfg.WorkAt = "08:30"
	cfg.LunchTimeHours = 1.5
	cfg.DailyWorkHours = 8.0
	cfg.Database.User = "erp"
	cfg.Database.Password = "lD9nAKQgYElV2tan"
	cfg.Database.Host = "127.0.0.1"
	cfg.Database.Database = "erp"
	cfg.Database.Charset = "utf8mb4"
	cfg.Database.ParseTime = "true"
	cfg.Database.Loc = "Local"
}

//GetEnv 取得環境參數
func GetEnv() *Config {
	return cfg
}
