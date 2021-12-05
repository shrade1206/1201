package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dsn struct {
	UserName     string
	PassWord     string
	Addr         string
	Port         int
	DataBase     string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

var DB *gorm.DB
var sqlDB *sql.DB

func InitMysql() (err error) {
	file, err := os.Open("./config/MySQL_Config.json")
	if err != nil {
		return
	}
	var dsn Dsn
	data := json.NewDecoder(file)
	err = data.Decode(&dsn)
	if err != nil {
		return
	}
	conn1 := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dsn.UserName, dsn.PassWord, dsn.Addr, dsn.Port)
	sqlDB, err := sql.Open("mysql", conn1)
	if err != nil {
		log.Printf("MySQL Connection Error : %s", err.Error())
		return
	}
	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + dsn.DataBase)
	if err != nil {
		log.Printf("CREATE DATABASE Error : %s", err.Error())
		return
	}

	// sqlDB.SetConnMaxLifetime(time.Duration(dsn.MaxLifetime) * time.Second)
	sqlDB.SetMaxOpenConns(dsn.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dsn.MaxIdleConns)

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dsn.UserName, dsn.PassWord, dsn.Addr, dsn.Port, dsn.DataBase)
	DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("MySQL Error : %s", err.Error())
		return
	}

	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Printf("AutoMigrate Error : %s", err.Error())
		return
	}
	return
}
