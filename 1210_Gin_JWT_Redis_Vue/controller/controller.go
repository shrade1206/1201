package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"token/redis"
	"token/util"

	"github.com/gin-contrib/sessions"
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

var UserData User
var Session sessions.Session

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
		check, err := redis.Client.Exists(reg.Username).Result()
		if err != nil {
			log.Printf("Param Error  :%s", err.Error())
			c.JSON(http.StatusOK, API_Error{
				Code: Code_Param_Invalid,
				Msg:  "Redis Exists Error :" + err.Error(),
			})
			return
		}
		// check = 0，資料不存在，check = 1，資料已存在
		if check == 0 {
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
					Code: Code_ok,
					Msg:  "Register Success",
				})
				log.Println("註冊成功")
			}
		} else {
			// 帳號已存在
			c.JSON(http.StatusOK, API_Error{
				Code: Code_Param_Invalid,
				Msg:  "UserName Already Exists",
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

		// 確認帳號是否存在
		// Exists查詢key是否存在，回傳true 或是false
		check, err := redis.Client.Exists(login.Username).Result()
		if err != nil {
			log.Printf("Redis EXists Error :%s", err.Error())
			c.JSON(http.StatusOK, API_Error{
				Code: Code_Param_Invalid,
				Msg:  "Redis EXists Error :" + err.Error(),
			})
			return
		}
		// check = 0，資料不存在，check = 1，資料已存在
		if check == 1 {
			value, err := redis.Client.Get(login.Username).Result()
			if err != nil {
				log.Printf("select Error :%s", err.Error())
				c.JSON(http.StatusOK, API_Error{
					Code: Code_DB_Conn,
					Msg:  "DB Error :" + err.Error(),
				})
				return
			} else {
				UserData = User{Username: login.Username, Password: value}
				// 資料庫取出的密碼 與 使用者輸入的密碼 驗證
				err := bcrypt.CompareHashAndPassword([]byte(UserData.Password), []byte(login.Password))
				if err != nil {
					c.JSON(http.StatusOK, API_Error{
						Code: Code_Param_Invalid,
						Msg:  "PassWord Error :" + err.Error(),
					})
					return
				} else {
					// 密碼對就生成Token回傳
					token, err := util.GenToken(UserData.Username)
					if err != nil {
						log.Printf("GenToken Error :%s" + err.Error())
						c.JSON(http.StatusOK, API_Error{
							Code: Code_Param_Invalid,
							Msg:  "Param Error :" + err.Error(),
						})
						return
					} else {
						// 登錄成功，生成session
						Session = sessions.Default(c)
						// 把username存進session
						Session.Set("username", UserData.Username)
						// 設定Session存活時間
						Session.Options(sessions.Options{
							MaxAge: 600,
						})
						err := Session.Save()
						if err != nil {
							log.Printf("Session Save Error :%s", err.Error())
							c.JSON(http.StatusOK, API_Error{
								Code: Code_Param_Invalid,
								Msg:  "Param Error :" + err.Error(),
							})
							return
						} else {
							redis.Client.Del("session_" + Session.ID())
						}
						// 把Token傳到前端cookie
						// cookie名字、cookie內容、存活時間，設定-1 = 刪除、path使用路經、作用host、是否只能https協議發送到服務端、HttpOnly=true 就不能被js獲取到、
						c.SetCookie("token", token, 7200, "/", "localhost", false, true)
						// 創好的Token當Key，session ID當Value存進Redis
						_, err = redis.Client.Set(token, Session.ID(), 0).Result()
						if err != nil {
							log.Printf("Save Token Error :%s", err.Error())
							c.JSON(http.StatusOK, API_Error{
								Code: Code_DB_Conn,
								Msg:  "DB Error :" + err.Error(),
							})
							return
						}
						c.JSON(http.StatusOK, API_Error{
							Code: Code_ok,
							Msg:  UserData.Username,
							Data: token,
						})
						log.Println(UserData.Username + " 登錄成功")
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

func Logout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusOK, API_Error{
			Code: Code_Param_Invalid,
			Msg:  "Cookie Error :" + err.Error(),
		})
		return
	}
	c.SetCookie("token", token, -1, "/", "localhost", false, true)
	redis.Client.Del("session_" + Session.ID())
	// redis.Client.Del(Session.ID())
	redis.Client.Del(token)
	c.JSON(200, API_Error{
		Code: Code_ok,
		Msg:  "Logout Success",
	})
}

func Noroute(c *gin.Context) {
	path := c.Request.URL.Path
	method := c.Request.Method
	fmt.Println(path)
	fmt.Println(method)
	//檢查path的開頭使是否為"/"
	if strings.HasPrefix(path, "/") {
		fmt.Println("Route ok")
	}
}
