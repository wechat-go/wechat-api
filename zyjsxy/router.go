package main

import (
	"encoding/base64"
	"net/http"
	"strings"
	"zyjsxy-api/controllers/login"
	"zyjsxy-api/controllers/personneloffice"
	"zyjsxy-api/controllers/system"
	"zyjsxy-api/controllers/workflow"
	"zyjsxy-api/controllers/workreport"
	"zyjsxy-api/database"
	"zyjsxy-api/util/aes"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

var (
	rd *redis.Pool
)

func init() {
	rd = database.Rd
}
func initRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())

	//用户登录、注册、密码找回
	r.POST("/api/oauth/token", login.GetTokenApi)

	//验证中间件
	r.Use(authMiddleWare())

	//用户退出
	r.POST("/api/revoke/token", login.RevokeTokenApi)

	return r
}
func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authStr := c.GetHeader("Authorization")
		if authStr != "" {
			aesEnc := aes.AesEncrypt{}
			tokenori, err0 := base64.StdEncoding.DecodeString(authStr)
			tokenstr, err := aesEnc.Decrypt(tokenori)
			if err != nil || err0 != nil {
				c.JSON(403, "权限受限，请重新登录！")
				c.Abort()
			} else {
				key := "token_" + handleStr(tokenstr, "id:", ",ip:")
				drtoken := queryrd(key)
				c.Set("requestId", handleStr(tokenstr, "id:", ",ip:"))
				if drtoken == authStr {
					c.Next()
				} else {
					c.JSON(403, "权限受限，请重新登录！")
					c.Abort()
				}
			}
		} else {
			c.JSON(403, "您无权访问，请登录！")
			c.Abort()
		}
	}
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("content-type", "application/json; charset=utf-8")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
			c.Abort()
		}
		c.Next()
	}
}

func handleStr(all string, target string, endstr string) string {
	if strings.Index(all, target) == -1 {
		return "-1"
	} else {
		start := strings.Index(all, target) + len(target)
		newStr := all[start:]
		end := strings.Index(newStr, endstr)
		return newStr[:end]
	}
}
func queryrd(t string) string {
	conn := rd.Get()
	defer conn.Close() //redis操作
	v, err := redis.String(conn.Do("GET", t))
	if err != nil {
		return ""
	}
	return v
}
