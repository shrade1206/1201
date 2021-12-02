package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// type Client struct {
// 	Id       int    `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }
type RegisterData struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginData struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// var users []User

func main() {
	err := InitMysql()
	if err != nil {
		log.Printf("initMysql() invalid : %s", err.Error())
		return
	}
	defer sqlDB.Close()

	r := gin.Default()
	// 註冊帳號--------------------------------
	r.POST("/Register", func(c *gin.Context) {
		var reg RegisterData
		var user User
		err := c.BindJSON(&reg)
		if err != nil {
			log.Printf("Register BindJson Error : %s", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		}
		model := GetUsername(reg.Username)
		if model.Id == 0 {
			user = User{Username: reg.Username, Password: reg.Password}
			err = DB.Create(&user).Error
			if err != nil {
				log.Printf("Register Error :%s", err.Error())
				return
			}
		} else {
			log.Println("帳戶已存在")
			return
		}
	})
	// 登錄帳號--------------------------------
	r.POST("/login", func(c *gin.Context) {
		var login LoginData
		err := c.BindJSON(&login)
		if err != nil {
			log.Printf("login BindJson Error : %s", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		}
		model := GetUsername(login.Username)
		if model.Id == 0 {
			log.Println("查無此帳號")
			return
		}
		if model.Password == login.Password {
			log.Println("登錄成功")
		} else {
			log.Println("密碼錯誤")
			return
		}
	})
	// --------------------------------
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return
	}
}

func GetUsername(username string) User {
	var user User
	err := DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		log.Printf("select Error :%s", err.Error())
	}
	return user
}
