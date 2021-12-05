package main

import (
	"fmt"
	"log"
	"net/http"
	"token/middleware"
	"token/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	Username string `json:"username" form:"username" validate:"required,min=6,max=12"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=12"`
}

type LoginData struct {
	Username string `json:"username" form:"username" validate:"required,min=6,max=12"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=12"`
}

type API_Error struct {
	Code int
	Msg  string
	Data interface{}
}

var validate = validator.New()

const Code_Param_Invalid = 1
const Code_DB_Conn = 2

// var users []User

func main() {
	err := InitMysql()
	if err != nil {
		log.Printf("initMysql() invalid : %s", err.Error())
		return
	}
	defer sqlDB.Close()

	r := gin.Default()
	// 測試生成--------------------------------
	r.GET("/test", func(c *gin.Context) {
		token, _ := util.GenToken("test")
		fmt.Println("Bearer " + token)
	})
	// 測試解析--------------------------------
	r.GET("/parse", func(c *gin.Context) {
		var tk util.TokenData
		err = c.BindJSON(&tk)
		if err != nil {
			log.Println("解析失敗")
			return
		}
		a, _ := util.ParseToken(tk.Token)
		fmt.Println(a.Username, a.Issuer, a.ExpiresAt)
	})
	// 測試認證01--------------------------------
	r.GET("/auth01", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		// 因使用c.next()，只要沒有關掉伺服器，註冊完即可使用，如果有關掉的話會報錯 Key "username" does not exist
		// username := c.MustGet("username").(string)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "success",
		})
	})
	// 測試認證02--------------------------------
	// r.GET("/auth02", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"code": 1,
	// 		"msg":  "success",
	// 	})
	// }).Use(middleware.JWTAuthMiddleware())

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
		// 驗證完struct格式是否正確，才繼續執行
		u := RegisterData{Username: reg.Username, Password: reg.Password}
		err = validate.Struct(u)
		if err != nil {
			log.Printf("validate Error :%s", err.Error())
			return
		} else {
			// 驗證帳號是否已存在，不存在就Create
			model := GetUsername(reg.Username)
			if model.Id == 0 {
				user = User{Username: reg.Username, Password: reg.Password}
				err = DB.Create(&user).Error
				if err != nil {
					log.Printf("Register Error :%s", err.Error())
					return
				}
			} else {
				log.Println("帳號已存在")
				return
			}
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
		u := LoginData{Username: login.Username, Password: login.Password}
		err = validate.Struct(u)
		if err != nil {
			log.Printf("validate Error :%s", err.Error())
			return
		} else {
			// 確認帳號是否存在
			model := GetUsername(login.Username)
			if model.Id == 0 {
				log.Println("查無此帳號")
				return
			}
			// 密碼對就生成Token回傳
			if model.Password == login.Password {
				token, err := util.GenToken(model.Username)
				if err != nil {
					log.Printf("GenToken Error :%s" + err.Error())
					return
				} else {
					// 把token傳到前端，保存在cookie
					c.JSON(http.StatusOK, gin.H{
						"username": model.Username,
						"token":    "Bearer " + token,
					})
					log.Println(model.Username + " 登錄成功")
				}
			} else {
				log.Println("密碼錯誤")
				return
			}
		}
	})
	// --------------------------------
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return
	}
}

// 用來資料校對
func GetUsername(username string) User {
	var user User
	err := DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		log.Printf("select Error :%s", err.Error())
		user = User{}
		return user
	}
	return user
}
