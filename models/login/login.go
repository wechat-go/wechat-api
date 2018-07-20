package login

import (
	"fmt"
	"strconv"
	db "zyjsxy-api/database"

	"github.com/garyburd/redigo/redis"
)

func GetToken(t int) string {
	conn := db.Rd.Get()
	defer conn.Close() //redis操作
	qr := "token_" + strconv.Itoa(t)
	v, err := redis.String(conn.Do("GET", qr))
	if err != nil {
		return ""
	}
	return v
}

func SaveToken(a, b string, c int) {
	conn := db.Rd.Get()
	defer conn.Close() //redis操作
	conn.Do("SETEX", a, c, b)
}

func GetTokenTTL(t string) int64 {
	conn := db.Rd.Get()
	defer conn.Close() //redis操作
	v, err := redis.Int64(conn.Do("TTL", t))
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return v
}

func DelToken(key string) error {
	conn := db.Rd.Get()
	defer conn.Close() //redis操作
	_, err := conn.Do("DEL", key)
	return err
}
