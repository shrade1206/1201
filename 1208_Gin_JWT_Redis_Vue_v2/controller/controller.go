package controller

import (
	"log"
	"net/http"
	"token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterData struct {
	Username string `json:"username" form:"username" binding:"required,min=6,max=12"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
}

type LoginData struct {
	Username string `json:"username" form:"username" binding:"required,min=6,max=12"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
}

type API_Error struct {
	Code int
	Msg  string
	Data interface{}
}

const Code_ok = 1
const Code_Param_Invalid = 2
const Code_DB_Conn = 3

func Register(c *gin.Context) {
	var reg RegisterData
	// 使用者傳帳密進來，解析到struct
	err := c.ShouldBindJSON(&reg)
	if err != nil {
		log.Printf("Register BindJson Error : %s", err.Error())
		c.JSON(http.StatusOK, API_Error{
			Code: Code_Param_Invalid,
			Msg:  "BindJSON Data:" + err.Error(),
		})
		return
	} else {
		// 驗證帳號是否已存在，不存在就Create
		model, err := redis.Client.Exists(reg.Username).Result()
		if err != nil {
			log.Printf("Param Error  :%s", err.Error())
			c.JSON(http.StatusOK, API_Error{
				Code: Code_Param_Invalid,
				Msg:  "Redis Exists Error :" + err.Error(),
			})
			return
		}
		// model = 0，資料不存在，model = 1，資料已存在
		if model == 0 {
			// 密碼加密
			hash, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Param Error  :%s", err.Error())
				c.JSON(http.StatusOK, API_Error{
					Code: Code_Param_Invalid,
					Msg:  "PassWord Encryption :" + err.Error(),
				})
				return
			}
			// 轉成string再保存
			regPassword := string(hash)
			// 創建帳密key、value
			SaveData, err := redis.Client.Set(reg.Username, regPassword, 0).Result()
			if err != nil {
				log.Printf("Register Error :%s", err.Error())
				c.JSON(http.StatusOK, API_Error{
					Code: Code_DB_Conn,
					Msg:  "Redis Set Error :" + err.Error(),
				})
				return
			}
			// 長度 > 0 = 保存成功
			if len(SaveData) > 0 {
				c.JSON(http.StatusOK, API_Error{
					Msg: "Register Success",
				})
			}
		} else {
			// 帳號已存在
			c.JSON(http.StatusOK, gin.H{
				"Msg": "UserName Already Exists",
			})
			return
		}
	}
}

func Login(c *gin.Context) {
	var login LoginData
	// 使用者傳帳密進來，解析到struct
	err := c.ShouldBindJSON(&login)
	if err != nil {
		log.Printf("login BindJson Error : %s", err.Error())
		c.JSON(http.StatusOK, API_Error{
			Code: Code_Param_Invalid,
			Msg:  "Param Error :" + err.Error(),
		})
		return
	} else {
		var user User
		// 確認帳號是否存在
		// Exists查詢key是否存在，回傳true 或是false
		sel, err := redis.Client.Exists(login.Username).Result()
		if err != nil {
			log.Printf("Redis EXists Error :%s", err.Error())
			c.JSON(http.StatusOK, API_Error{
				Code: Code_Param_Invalid,
				Msg:  "Redis EXists Error :" + err.Error(),
			})
			return
		}
		if sel == 1 {
			value, err := redis.Client.Get(login.Username).Result()
			if err != nil {
				log.Printf("select Error :%s", err.Error())
				c.JSON(http.StatusOK, API_Error{
					Code: Code_DB_Conn,
					Msg:  "DB Error :" + err.Error(),
				})
				return
			} else {
				user = User{Username: login.Username, Password: value}
				// 資料庫取出的密碼 與 使用者輸入的密碼 驗證
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
				if err != nil {
					c.JSON(http.StatusOK, API_Error{
						Code: Code_Param_Invalid,
						Msg:  "PassWord Error :" + err.Error(),
					})
					return
				} else {
					// 密碼對就生成Token回傳
					token, err := util.GenToken(user.Username)
					if err != nil {
						log.Printf("GenToken Error :%s" + err.Error())
						c.JSON(http.StatusOK, API_Error{
							Code: Code_Param_Invalid,
							Msg:  "Param Error :" + err.Error(),
						})
						return
					} else {
						// 把Token傳到前端cookie
						// cookie名字、cookie內容、存活時間，設定-1 = 刪除、path使用路經、作用host、是否只能https協議發送到服務端、HttpOnly=true 就不能被js獲取到、
						c.SetCookie("token", token, 7200, "/", "localhost", false, true)
						// 創好的Token存進Redis
						_, err := redis.Client.Set(user.Username+"Token", token, 0).Result()
						if err != nil {
							log.Printf("Save Token Error :%s", err.Error())
							c.JSON(http.StatusOK, API_Error{
								Code: Code_DB_Conn,
								Msg:  "DB Error :" + err.Error(),
							})
							return
						}
						// 把token傳到前端，保存在cookie
						c.JSON(http.StatusOK, API_Error{
							Msg: user.Username + " 登錄成功",
							// "username": model.Username,
							Data: "Bearer " + token,
						})
						log.Println(user.Username + " 登錄成功")
					}
				}
			}
		} else {
			log.Printf("查無此帳號")
			c.JSON(http.StatusOK, API_Error{
				Code: Code_Param_Invalid,
				Msg:  "查無此帳號",
			})
			return
		}

	}
}
