package db

import (
	"fmt"

	"github.com/lazyironf4ur/go-infra/common"
	"github.com/lazyironf4ur/go-infra/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlDB *gorm.DB

func init() {
	if conf.GlobalConfig["mysql"] != nil {
		mysqlConf := conf.GlobalConfig["mysql"].(map[string]interface{})
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlConf["username"], mysqlConf["password"], mysqlConf["host"], mysqlConf["port"], mysqlConf["database"])
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: glog,
		})
		if err != nil {
			panic(err)
		}
		MysqlDB = db
		common.Must(MysqlDB)
	}
}
