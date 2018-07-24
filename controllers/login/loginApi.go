package login

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"
	. "zyjsxy-api/models/login"
	"zyjsxy-api/models/system"
	"zyjsxy-api/util"
	"zyjsxy-api/util/aes"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func GetTokenApi(c *gin.Context) {
	obj := util.GetObj(c.Request.Body)
	sort := int(obj["sort"].(float64))
	if sort == 1 {
		register(c, &obj)
	} else if srot == 2 {
		login(c, &obj)
	} else if sort == 3 {
		resetpw(c, &obj)
	} else if sort == 4 || sort == 5 { //注册短信发送
		sendSms(c, obj)
	}
}

func RevokeTokenApi(c *gin.Context) {
	authStr := c.GetHeader("Authorization")
	if authStr != "" {
		aesEnc := aes.AesEncrypt{}
		tokenori, err0 := base64.StdEncoding.DecodeString(authStr)
		tokenstr, err := aesEnc.Decrypt(tokenori)
		if err != nil || err0 != nil {
			c.JSON(403, "权限受限，请重新登陆！")
		} else {
			key := "token_" + util.HandleStr(tokenstr, "id:", ",ip:")
			err = DelToken(key)
			if err == nil {
				c.JSON(200, "退出成功！")
			} else {
				c.JSON(403, "退出失败！")
			}

		}
	} else {
		c.JSON(403, "权限受限，请重新登陆！")
	}
}

func sendSms(c *gin.Context, o *map[string]interface{}) {
	//参数获取
	obj := util.GetObj(c.Request.Body)
	sort := int(obj["sort"].(float64))
	phone := obj["phone"].(string)

	//开始发送
	if !util.IsPhone(phone) {
		c.JSON(500, "手机号码不对，请重新输入！")
		return
	}

	//获取手机号码发送的次数
	np := GetRedis("np_" + phone)
	if np == "3" {
		c.JSON(500, "每天最多只能发送三条短信，您已超过次数！")
		return
	}

	info := util.NewSmsInfo()
	info.Mobile = phone
	if sort == 4 {
		info.Sort = "注册"
	} else if sort == 5 {
		info.Sort = "密码重置"
	}
	b := uitl.Sendsms(info)
	if b {
		c.JSON(200, "验证码发送成功！")
		SaveRedis("nr_"+phone, info.Code, 300)
		if np == "" {
			SaveRedis("np_"+phone, "1", 86400)
		} else if np == "1" {
			UpRedis("np_"+phone, "2")
		} else if np == "2" {
			UpRedis("np_"+phone, "3")
		}
		return
	} else {
		c.JSON(500, "您的手机号暂时不支持！")
		return
	}
}

func register(c *gin.Context, o *map[string]interface{}) {
	phone := obj["phone"].(string)
	pass := obj["pass"].(string)

	//判断是否非空
	if phone == "" || pass == "" {
		c.JSON(500, "手机号或密码不能为空！")
		return
	}

	//新建用户
	user := system.User{Username: phone, phone: phone, Avatar: "/img/user.jpg", Password: aes.Hashsalt(pass)}
	err := system.AddUser(&user)
	if err != nil {
		c.JSON(500, "您的注册审核不通过，请联系管理员！")
		return
	}

	//结果构造
	results := make(map[string]interface{})

	//获取用户ip
	orip := c.GetHeader("X-Real-IP")
	if orip == "" {
		orip = c.Request.RemoteAddr
	}

	//加密包
	aesEnc := aes.AesEncrypt{}

	//token生成
	str := "{id:" + strconv.Itoa(int(user.ID)) + ",ip:" + orip + ",over:" + strconv.FormatInt(time.Now().Unix()+7200, 10) + "}"
	tokenbyte, _ := aesEnc.Encrypt(str)
	toke := base64.StdEncoding.EncodeToString(tokenbyte)
	results["access_token"] = toke
	results["expires_in"] = 7200
	SaveToken("token_"+strconv.Itoa(int(u.ID)), toke, 7200)

	c.JSON(200, results)

}

func login(c *gin.Context, o *map[string]interface{}) {
	userName := obj["username"].(string)
	passWord := obj["password"].(string)
	if userName != "" && passWord != "" {
		user := system.User{Username: userName}
		u, err := system.GetUserFrist(user)
		if err != nil {
			if err.Error() == "record not found" {
				c.JSON(500, "用户不存在！")
			} else {
				c.JSON(500, "服务器错误！")
			}
		} else {
			cop := aes.Compare(u.Password, passWord)
			if cop == nil {
				drtoken := GetToken(int(u.ID))
				results := make(map[string]interface{})
				orip := c.GetHeader("X-Real-IP")
				if orip == "" {
					orip = c.Request.RemoteAddr
				}
				if strings.Index(orip, ":") != -1 {
					orip = orip[:strings.Index(orip, ":")]
				}
				var tokenip string
				aesEnc := aes.AesEncrypt{}
				if drtoken != "" {
					tokenori, _ := base64.StdEncoding.DecodeString(drtoken)
					tokenstr, _ := aesEnc.Decrypt(tokenori)
					tokenip = util.HandleStr(tokenstr, "ip:", ",")
				}
				if drtoken == "" || orip != tokenip {
					str := "{id:" + strconv.Itoa(int(u.ID)) + ",ip:" + orip + ",over:" + strconv.FormatInt(time.Now().Unix()+7200, 10) + "}"
					tokenbyte, _ := aesEnc.Encrypt(str)
					toke := base64.StdEncoding.EncodeToString(tokenbyte)
					results["access_token"] = toke
					results["expires_in"] = 7200
					SaveToken("token_"+strconv.Itoa(int(u.ID)), toke, 7200)
				} else {
					ti := GetTokenTTL("token_" + strconv.Itoa(int(u.ID)))
					results["access_token"] = drtoken
					results["expires_in"] = ti
				}
				c.JSON(200, results)
			} else {
				c.JSON(400, "密码错误！")
			}
		}
	} else {
		c.JSON(400, "用户名或密码不能为空！")
	}
}

func resetpw(c *gin.Context, o *map[string]interface{}) {

}
