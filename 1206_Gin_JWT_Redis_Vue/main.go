package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"token/middleware"

	// "token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterData struct {
	Username string `json:"username" form:"username" validate:"required,min=6,max=12"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=12"`
}

type LoginData struct {
	Username string `json:"username" form:"username" validate:"required,min=6,max=12"`
	Password string `json:"password" form:"password" validate:"required,min=6,max=12"`
}

// type API_Error struct {
// 	Code int
// 	Msg  string
// 	Data interface{}
// }

var Client *redis.Client
var validate = validator.New()

// const Code_Param_Invalid = 1
// const Code_DB_Conn = 2

func main() {
	Client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //不設定密碼
		DB:       0,  //預設資料庫
	})
	// err := redis.InitRedis()
	// if err != nil {
	// 	log.Printf("InitRedis Error :%s", err.Error())
	// 	return
	// }
	// defer redis.Client.Close()

	r := gin.Default()
	// 測試------------------------------------
	r.GET("testdel", func(c *gin.Context) {
		Client.Del("jiyuusama")
	})
	r.POST("test1", func(c *gin.Context) {
		var reg RegisterData
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
		}
		a, _ := Client.Set(reg.Username, reg.Password, 0).Result()
		b := Client.Get(reg.Username)
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
	// 測試認證01--------------------------------
	r.GET("/auth01", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		// 因使用c.next()，只要沒有關掉伺服器，註冊完即可使用，如果有關掉的話會報錯 Key "username" does not exist
		// username := c.MustGet("username").(string)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "success",
		})
	})
	// 註冊帳號--------------------------------
	r.POST("/Register", func(c *gin.Context) {
		var reg RegisterData
		err := c.BindJSON(&reg)
		if err != nil {
			log.Printf("Register BindJson Error : %s", err.Error())
			c.JSON(http.StatusOK, err.Error())
			return
		}
		u := RegisterData{Username: reg.Username, Password: reg.Password}
		err = validate.Struct(u) // 確認 struct 格式是否正確
		if err != nil {
			log.Printf("validate Error :%s", err.Error())
			return
		} else {
			// 驗證帳號是否已存在，不存在就Create
			model, err := Client.Exists(reg.Username).Result()
			if err != nil {
				log.Printf("參數錯誤 :%s", err.Error())
				return
			}
			if model == 0 {
				//創建帳密key、value
				r, err := Client.Set(reg.Username, reg.Password, 0).Result()
				if len(r) > 0 {
					log.Println("註冊成功")
				}
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
		err = validate.Struct(u) // 確認 struct 格式是否正確
		if err != nil {
			log.Printf("validate Error :%s", err.Error())
			return
		} else {
			// 確認帳號是否存在
			model := GetUsername(login.Username)
			// 密碼對就生成Token回傳
			if model.Password == login.Password {
				token, err := util.GenToken(model.Username)
				if err != nil {
					log.Printf("GenToken Error :%s" + err.Error())
					return
				} else {
					// 創好的Token存進Redis
					_, err := Client.Set(model.Username+"Token", token, 0).Result()
					if err != nil {
						log.Printf("Save Token Error :%s", err.Error())
						return
					}
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
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("8080 err : ", err.Error())
		return
	}
}

// 取得資料
func GetUsername(username string) User {
	Client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //不設定密碼
		DB:       0,  //預設資料庫
	})
	var user User
	// Exists查詢key是否存在，回傳true 或是false
	sel, err := Client.Exists(username).Result()
	if err != nil {
		log.Printf("參數錯誤 Error :%s", err.Error())
		return user
	}
	if sel == 1 {
		value, err := Client.Get(username).Result()
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
