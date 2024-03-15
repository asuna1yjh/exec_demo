package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var declareUrl = "https://poke.migu.cn/poke/media/upload/v1.0"

type Resp struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

// Declare 申报文件上传的方法
func Declare(cookie string, parameter string) error {
	url := "https://poke.migu.cn/poke/vrbt/video/save"
	method := "POST"

	payload := strings.NewReader(parameter)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("pacmtoken", cookie)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var resp = new(Resp)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return err
	}
	if resp.Code != "000000" {
		return fmt.Errorf("declare failed, code:%v, info:%v", resp.Code, resp.Info)
	}
	return err
}
