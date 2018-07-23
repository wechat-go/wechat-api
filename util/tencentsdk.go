// tencentsdk
package util

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
)

var (
	appid     = "1400112068"
	appscript = "eb1b91d9f8547c6e20574f8b10d53113"
)

type SmsInfo struct {
	Sort   string //类型：登录、注册、密码找回
	Code   string //验证码内容
	Pri    string //有效时间：5
	Mobile string //手机号
}

func NewSmsInfo() *SmsInfo {
	return &SmsInfo{
		Sort: "",
		Code: strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)),
		Pri:  strconv.Itoa(5),
	}
}

func Sendsms(info *SmsInfo) bool { // 指定模板发送短信
	unix := time.Now().Unix()

	s := "appkey=" + appscript + "&random=" + info.Code + "&time=" + strconv.Itoa(unix) + "&mobile=" + info.Mobile

	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	request := gorequest.New()

	resp, body, errs := request.Post("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=" + appid + "&random=" + info.Code).
		Send(`{"ext":"",
		 	"extend":"",
		    "params": [
		        "` + info.Sort + `",
		        "` + info.Code + `",
		        "` + info.Pri + `"
		    ],
			"sig":"` + string(bs) + `",
		    "tel": {
		        "mobile": "` + info.Mobile + `",
		        "nationcode": "86"
		    },
		    "time": ` + strconv.Itoa(time.Unix()) + `,
		    "tpl_id": 159496
		}`).
		End()
	var obj map[string]interface{}
	json.Unmarshal([]byte(body), &obj)
	if int(obj["result"].(float64)) == 0 && obj["errmsg"].(string) == "OK" {
		return true
	}
	return false
}
