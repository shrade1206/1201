package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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
var SQLDB *sql.DB

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
	SQLDB, err = sql.Open("mysql", conn1)
	if err != nil {
		return
	}
	_, err = SQLDB.Exec("CREATE DATABASE IF NOT EXISTS " + dsn.DataBase)
	if err != nil {
		return
	}

	SQLDB.SetConnMaxLifetime(time.Duration(dsn.MaxLifetime) * time.Second)
	SQLDB.SetMaxOpenConns(dsn.MaxOpenConns)
	SQLDB.SetMaxIdleConns(dsn.MaxIdleConns)

	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dsn.UserName, dsn.PassWord, dsn.Addr, dsn.Port, dsn.DataBase)
	DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return
	}
	err = DB.AutoMigrate(&User{})
	if err != nil {
		return
	}
	return
}
