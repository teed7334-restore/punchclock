package beans

//Holiday 人事行政局回傳資料格式
type Holiday struct {
	Date        string `json:"Start Date"`
	Name        string `json:"Subject"`
	Description string `json:"Description"`
}
