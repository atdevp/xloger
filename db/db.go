package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/xloger/g"
)

type DBPool struct {
	Xlog *gorm.DB
}

var (
	dbp DBPool
)

func Con() DBPool {
	return dbp
}

func InitDB() {
	xlog, err := gorm.Open("mysql", g.Config().MysqlConfig.Dsn)
	// xlog.LogMode(true)       // 显示详细日志
	xlog.SingularTable(true) //  如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	xlog.DB().SetMaxOpenConns(20)

	if err != nil {
		log.Fatalln("con mysql err: ", err)
	}
	dbp.Xlog = xlog
	log.Printf("conntected mysqld")
}
