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

func register(c *gin.Context, o *map[string]interface{}) {
	phone := o["phone"].(string)
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
