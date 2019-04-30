package db

import (
	"fmt"
	"log"
	"time"

	"github.com/Bo0km4n/claude/app/common/proto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Mysql *gorm.DB

func InitMysql() {
	dialect := "mysql"
	host := "127.0.0.1"
	port := "3306"
	user := "claude"
	password := "password"
	database := "claude"
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, password, host, port, database)

	db, err := gorm.Open(dialect, url)
	if err != nil {
		log.Fatal("MYSQL ERROR: ", err)
	}
	db.SingularTable(true)
	db.BlockGlobalUpdate(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetConnMaxLifetime(time.Duration(60) * time.Second)
	Mysql = db
}

func MigrateMysql() {
	Mysql.AutoMigrate(
		&proto.LREntry{},
	)
}

func TestInitMysql() {
	dialect := "mysql"
	host := "127.0.0.1"
	port := "3306"
	user := "root"
	password := "password"
	database := "claude_test"
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, password, host, port, database)

	db, err := gorm.Open(dialect, url)
	if err != nil {
		log.Fatal("MYSQL ERROR: ", err)
	}
	db.SingularTable(true)
	db.BlockGlobalUpdate(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetConnMaxLifetime(time.Duration(60) * time.Second)
	Mysql = db
}

func TestCloseMysql() {
	Mysql.DropTable("lr_entry")
}
