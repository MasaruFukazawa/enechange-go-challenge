package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	config "go-challenge/config"
)

var (
	dbMy *gorm.DB
	err  error
)

func Init() {
	openMySQL()
}

func Close() {
	my, _ := dbMy.DB()
	my.Close()
}

// GetDB はデータベース接続を返す
func GetDB() *gorm.DB {
	return dbMy
}

func openMySQL() {
	cnf := config.GetConfig()

	database := cnf.MySQLDatabase

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cnf.MySQLUser,
		cnf.MySQLPassword,
		cnf.MySQLHost,
		cnf.MySQLPort,
		database,
	)

	dbMy, err = gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic(err)
	}
}
