package request

import (
	"fmt"
	"testing"
)

func TestDeclare(t *testing.T) {
	var a = `{
    "publishTemplateId": "95054709-6b60-437f-8222-dd9a7f362219",
    "videoFileInfo": {
        "videoName": "测试足球栏目上传03",
        "filePath": "/video/2024/03/13/15/10/14615232-917f-48c8-a34e-f8c22dd1e8af.mp4",
        "duration": "264"
    }
}`
	err := Declare("C9948B8891A2A58A69958DA2867993729A94848994A2A98C6896899D8075A0729C97848F97A9A688668D8D9D888197718FC9B5C3D4D5-3974503462", a)
	if err != nil {
		t.Logf("err: %v", err)
		return
	}
	fmt.Println("success")

}
