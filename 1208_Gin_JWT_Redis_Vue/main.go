package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

const Code_Param_Invalid = 1
const Code_DB_Conn = 2
const Code_Server_Conn = 3

func main() {
	err := redis.InitRedis()
	if err != nil {
		log.Printf("InitRedis Error :%s", err.Error())
		return
	}
	defer redis.Client.Close()

	r := gin.Default()
	// 測試------------------------------------
	r.GET("ttt", func(c *gin.Context) {
		u0 := User{}
		u0.Password = "qweqweqwe"
		hash, err := bcrypt.GenerateFromPassword([]byte(u0.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		pwd := string(hash)
		fmt.Println(pwd)

		u1 := User{}
		u1.Password = pwd
		login := "qweqweqwe"
		// 資料庫取出的密碼 與 使用者輸入的密碼 驗證
		err = bcrypt.CompareHashAndPassword([]byte(u1.Password), []byte(login))
		if err != nil {
			log.Println(err.Error())

			return
		} else {
			log.Println("ok")
		}

	})
	r.GET("ttt1", func(c *gin.Context) {
		data, _ := c.Cookie("token")
		Data, err := util.ParseToken(data)
		if err != nil {
			return
		} else {
			fmt.Println(Data)
			c.SetCookie("token", data, -1, "/", "127.0.0.1", false, true)
		}
	})
	r.POST("test1", func(c *gin.Context) {
		var reg RegisterData
		err := c.ShouldBindJSON(&reg)
		if err != nil {
			log.Printf("Register BindJson Error : %s", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		}
		a, _ := redis.Client.Set(reg.Username, reg.Password, 0).Result()
		b := redis.Client.Get(reg.Username)
		fmt.Println(len(a))
		fmt.Println(b)
	})
	// 測試生成--------------------------------
	r.GET("/test", func(c *gin.Context) {
		token, _ := util.GenToken("test")
		fmt.Println("Bearer " + token)
	})
	// 測試解析--------------------------------
	r.GET("/parse", func(c *gin.Context) {
		var tk util.TokenData
		err := c.BindJSON(&tk)
		if err != nil {
			log.Println("解析失敗")
			return
		}
		a, _ := util.ParseToken(tk.Token)
		fmt.Println(a.Username, a.Issuer, a.ExpiresAt)
	})
	// 測試認證--------------------------------
	checkUser := func(c *gin.Context) (string, bool) {
		var ok bool
		u, ok := c.Get("username")
		if !ok {
			return "", false
		}
		var username string
		username, ok = u.(string) // panic
		if !ok {
			return "", false
		}
		return username, true
	}
	handler := func(c *gin.Context) {
		// 因使用c.next()，只要沒有關掉伺服器，註冊完即可使用，如果有關掉的話會報錯 Key "username" does not exist
		var username, ok = checkUser(c)
		if !ok {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "not login",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "success, " + username,
		})
	}
	r.GET("/auth01", handler)
	r.GET("/auth02", handler)
	// 註冊帳號--------------------------------
	r.POST("/Register", func(c *gin.Context) {
		var reg RegisterData
		err := c.ShouldBindJSON(&reg)
		if err != nil {
			log.Printf("Register BindJson Error : %s", err.Error())
			c.JSON(http.StatusOK, API_Error{
				Code: Code_Param_Invalid,
				Msg:  "Param Error :" + err.Error(),
			})
			return
		} else {
			// 驗證帳號是否已存在，不存在就Create
			model, err := redis.Client.Exists(reg.Username).Result()
			if err != nil {
				log.Printf("Param Error  :%s", err.Error())
				c.JSON(http.StatusOK, API_Error{
					Code: Code_Param_Invalid,
					Msg:  "Param Error :" + err.Error(),
				})
				return
			}
			if model == 0 {
				// 密碼加密
				hash, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
				if err != nil {
					log.Printf("Param Error  :%s", err.Error())
					c.JSON(http.StatusOK, API_Error{
						Code: Code_Param_Invalid,
						Msg:  "Param Error :" + err.Error(),
					})
					return
				}
				regPassword := string(hash)
				//創建帳密key、value
				r, err := redis.Client.Set(reg.Username, regPassword, 0).Result()
				if err != nil {
					log.Printf("Register Error :%s", err.Error())
					c.JSON(http.StatusOK, API_Error{
						Code: Code_DB_Conn,
						Msg:  "DB set Error :" + err.Error(),
					})
					return
				}
				if len(r) > 0 {
					c.JSON(http.StatusOK, API_Error{
						Msg: "Register success",
					})
				}
			} else {
				c.JSON(http.StatusOK, gin.H{
					"Msg": "帳號已存在",
				})
				return
			}
		}
	})
	// 登錄帳號--------------------------------
	r.POST("/login", func(c *gin.Context) {
		var login LoginData
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
			model := GetUsername(login.Username)
			// 資料庫取出的密碼 與 使用者輸入的密碼 驗證
			err = bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(login.Password))
			if err != nil {
				c.JSON(http.StatusOK, API_Error{
					Code: Code_Param_Invalid,
					Msg:  "PassWord Error :" + err.Error(),
				})
				return
			} else {
				// 密碼對就生成Token回傳
				token, err := util.GenToken(model.Username)
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
					c.SetCookie("token", token, 7200, "/", "127.0.0.1", false, true)
					// 創好的Token存進Redis
					_, err := redis.Client.Set(model.Username+"Token", token, 0).Result()
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
						Msg: model.Username + " 登錄成功",
						// "username": model.Username,
						Data: "Bearer " + token,
					})
					log.Println(model.Username + " 登錄成功")
				}
			}
		}

	})
	// --------------------------------
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./public"))), func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		fmt.Println(path)
		fmt.Println(method)
		//檢查path的開頭使是否為"/"
		if strings.HasPrefix(path, "/") {
			fmt.Println("Route ok")
		}
	})
	// --------------------------------
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return
	}
}

// 取得資料
func GetUsername(username string) User {
	var user User
	// Exists查詢key是否存在，回傳true 或是false
	sel, err := redis.Client.Exists(username).Result()
	if err != nil {
		log.Printf("參數錯誤 Error :%s", err.Error())
		return user
	}
	if sel == 1 {
		value, err := redis.Client.Get(username).Result()
		if err != nil {
			log.Printf("select Error :%s", err.Error())
			return user
		} else {
			user = User{Username: username, Password: value}
			return user
		}
	} else {
		log.Printf("查無此帳號")

	}
	return user
}
