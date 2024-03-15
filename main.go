package main

import (
	"bufio"
	"excel_demo/pkg"
	request "excel_demo/pkg/requests"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Classification string // 分类
	UID            string // uid
	VisualID       string // 视频id
	Account        string // 账号
	Token          string // token
	Counter        int
	Works          Works
	Data           [][]string
}

type Works struct {
	Basket    int // 篮球作品数量
	Foot      int // 足球作品数量
	Billiards int // 台球作品数量
	Bolley    int // 排球作品数量
	Tennis    int // 网球作品数量
	Esports   int // 电子竞技作品数量
}

// 重试队列
type retryQueue struct {
	Param       []byte // 参数
	Token       string // token
	MediaAssets string
	Num         int
}

func main() {
	f, err := excelize.OpenFile("马甲账号内容规划.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	var data = make([]*User, 0, 20)
	for i := range rows {
		var user = &User{
			Classification: getCell(f, "B", i+1),
			UID:            getCell(f, "C", i+1),
			VisualID:       getCell(f, "D", i+1),
			Account:        getCell(f, "E", i+1),
			Counter:        toInt(getCell(f, "F", i+1)),
			Works: Works{
				Basket:    toInt(getCell(f, "G", i+1)),
				Foot:      toInt(getCell(f, "H", i+1)),
				Billiards: toInt(getCell(f, "I", i+1)),
				Bolley:    toInt(getCell(f, "J", i+1)),
				Tennis:    toInt(getCell(f, "K", i+1)),
				Esports:   toInt(getCell(f, "L", i+1)),
			},
			Token: getCell(f, "N", i+1),
		}
		if user.Classification == "" {
			continue
		}
		data = append(data, user)
	}
	openFile, err := os.Open("file_1.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(openFile)
	var Basket [][]string    // 篮球作品数量
	var Foot [][]string      // 足球作品数量
	var Billiards [][]string // 台球作品数量
	var Bolley [][]string    // 排球作品数量
	var Tennis [][]string    // 网球作品数量
	var Esports [][]string   // 电竞作品数量
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		line = strings.TrimRight(line, "\n")
		l := strings.Split(line, "|")
		switch l[0] {
		case "篮球":
			Basket = append(Basket, l)
		case "足球":
			Foot = append(Foot, l)
		case "台球":
			Billiards = append(Billiards, l)
		case "排球":
			Bolley = append(Bolley, l)
		case "电竞":
			Esports = append(Esports, l)
		case "网球":
			Tennis = append(Tennis, l)
		}
	}
	// 输出
	for _, v := range data {
		if v.Works.Basket > 0 {
			for _, v1 := range pkg.GetSlice(&Basket, v.Works.Basket) {
				v.Data = append(v.Data, v1)
			}
		}
		if v.Works.Foot > 0 {
			for _, v1 := range pkg.GetSlice(&Foot, v.Works.Foot) {
				v.Data = append(v.Data, v1)
			}
		}
		if v.Works.Billiards > 0 {
			for _, v1 := range pkg.GetSlice(&Billiards, v.Works.Billiards) {
				v.Data = append(v.Data, v1)
			}
			if v.Works.Bolley > 0 {
				for _, v1 := range pkg.GetSlice(&Bolley, v.Works.Bolley) {
					v.Data = append(v.Data, v1)
				}
			}
		}
		if v.Works.Tennis > 0 {
			for _, v1 := range pkg.GetSlice(&Tennis, v.Works.Tennis) {
				v.Data = append(v.Data, v1)
			}
		}
		if v.Works.Esports > 0 {
			for _, v1 := range pkg.GetSlice(&Esports, v.Works.Esports) {
				v.Data = append(v.Data, v1)
			}
		}
	}
	//for _, v := range data {
	//	fmt.Printf("account: %s|%v|%v\n", v.Account, v.Counter, len(v.Data))
	//}
	f1 := excelize.NewFile()
	// 创建一个工作表
	_, err = f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	failedFile, err := os.OpenFile("file_err.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer failedFile.Close()
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	var retry = make(chan *retryQueue, 1000)
	start_row := 2
	for _, v := range data {
		for _, v1 := range v.Data {
			//i, _ := strconv.ParseFloat(, 64)
			//i = i / 1000
			p := Param(v1[2], v1[5], v1[6])
			err = request.Declare(v.Token, string(p))
			if err != nil {
				//if v1[1] == "1055318824" {
				//failedFile.WriteString(fmt.Sprintf("%s|%s\n", v1[0], v1[5]))
				log.Printf("申报失败，视频id: %s\n", v1[2])
				retry <- &retryQueue{
					Param:       p,
					Token:       v.Token,
					MediaAssets: v1[2],
				}
			} else {
				fmt.Println(v1[2], "success")
			}
			// 写入excel文件逻辑
			v1 = append(v1, v.Classification, v.UID, v.VisualID, v.Account, v.Token)
			f1.SetSheetRow("Sheet1", fmt.Sprintf("%s%d", "A", start_row), &v1)
			start_row++
		}
	}
	if err := f1.SaveAs("文件内容规划详情.xlsx"); err != nil {
		fmt.Println(err)
		return
	}
label:
	for {
		select {
		case v := <-retry:
			if v.Num < 3 {
				err = request.Declare(v.Token, string(v.Param))
				if err != nil {
					//if v.MediaAssets == "23/24赛季WCBA常规赛第1轮全场集锦：天津冠岚75:91石家庄英励" {
					log.Printf("%v加入重试队列\n", v.MediaAssets)
					fmt.Println(string(v.Param), v.Token)
					retry <- &retryQueue{
						Param:       v.Param,
						Token:       v.Token,
						MediaAssets: v.MediaAssets,
						Num:         v.Num + 1,
					}
				}
			} else {
				log.Printf("重试次数超过3次，视频id: %s\n", v.MediaAssets)
				failedFile.WriteString(fmt.Sprintf("%v|%v\n", string(v.Param), v.Token))
			}
		default:
			if len(retry) == 0 {
				break label
			}
		}
	}
}

func Param(videoName, filePath, duration string) []byte {
	b := fmt.Sprintf(`{
	"publishTemplateId": "95054709-6b60-437f-8222-dd9a7f362219",
	"videoFileInfo": {
	   "videoName": "%s",
	   "filePath": "%s",
	   "duration": "%s"
	}}`, videoName, filePath, duration)
	return []byte(b)
}

func getCell(f *excelize.File, column string, i int) string {
	v, _ := f.GetCellValue("Sheet1", fmt.Sprintf("%s%d", column, i+1))
	return v
}

func toInt(d string) int {
	s, _ := strconv.Atoi(d)
	return s
}
