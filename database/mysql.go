package database

import (
	"log"

	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Orm    *gorm.DB
	config tomlConfig
	Rd     *redis.Pool
)

type tomlConfig struct {
	DB database `toml:"database"`
}

type database struct {
	Server   string
	Username string
	Password string
	Dataname string
	Rdserver string
	RdPW     string
}

func init() {
	var err error

	if _, err = toml.DecodeFile("./conf.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "w_" + defaultTableName
	}
	Orm, err = gorm.Open("mysql", config.DB.Username+":"+config.DB.Password+"@tcp("+config.DB.Server+")/"+config.DB.Dataname+"?charset=utf8&parseTime=True&loc=Local")
	//	defer Orm.Close()
	if err != nil {
		log.Fatalln(err)
	}
	Rd = NewPool(config.DB.Rdserver, config.DB.RdPW)
	Orm.DB().SetMaxIdleConns(10)
	Orm.DB().SetMaxOpenConns(100)

	if err := Orm.DB().Ping(); err != nil {
		log.Fatalln(err)
	}
}
