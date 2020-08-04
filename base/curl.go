package base

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

//GetURL 透過HTTP GET取得網頁資料
func GetURL(url string) []byte {
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

//PostURL 透過HTTP POST扔資料給特定網頁
func PostURL(url string, params []byte) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(params))
	request.Header.Set("X-Custom-Header", "counter")
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
