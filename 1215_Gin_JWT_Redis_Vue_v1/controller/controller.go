package controller

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

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

const Code_ok = 1
const Code_Param_Invalid = 2
const Code_DB_Conn = 3
const Code_Session_Invalid = 4

var UserData User

func ErrMsg(c *gin.Context, code int, msg string, data interface{}, err error) {
	c.AbortWithStatusJSON(http.StatusOK, API_Error{
		Code: code,
		Msg:  msg + " " + err.Error(),
		Data: data,
	})
}
func Msg(c *gin.Context, code int, msg string, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, API_Error{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func OkMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, API_Error{
		Code: Code_ok,
		Msg:  msg,
		Data: data,
	})
}

func RedisExists(c *gin.Context, username string) (int64, error) {
	check, err := redis.Client.Exists(username).Result()
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "Redis", nil, err)
		return check, err
	}
	return check, nil
}
func RedisSet(c *gin.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	SaveData, err := redis.Client.Set(key, value, expiration).Result()
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "Redis", nil, err)
		return SaveData, err
	}
	return SaveData, nil
}

func Register(c *gin.Context) {
	var reg RegisterData
	// 使用者傳帳密進來，解析到struct
	err := c.ShouldBindJSON(&reg)
	if err != nil {
		ErrMsg(c, Code_Param_Invalid, "帳號密碼格式錯誤", nil, err)
		return
	}
	// 驗證帳號是否已存在，不存在就Create
	// check, err := RedisExists(c, reg.Username)
	check, err := redis.Client.Exists(reg.Username).Result()
	if err != nil {
		ErrMsg(c, Code_Param_Invalid, "Redis Exists", nil, err)
		return
	}
	// check = 0，資料不存在，check = 1，資料已存在
	if check == 0 {
		// 密碼加密
		hash, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
		if err != nil {
			ErrMsg(c, Code_Param_Invalid, "資料錯誤", nil, err)
			return
		}
		// 轉成string再保存
		regPassword := string(hash)
		// 創建帳密key、value
		SaveData, err := RedisSet(c, reg.Username, regPassword, 0)
		// SaveData, err := redis.Client.Set(reg.Username, regPassword, 0).Result()
		if err != nil {
			ErrMsg(c, Code_DB_Conn, "Redis", nil, err)
			return
		}
		// 長度 > 0 = 保存成功
		if len(SaveData) > 0 {
			OkMsg(c, "註冊成功", nil)
		}
	} else {
		// 帳號已存在
		Msg(c, Code_Param_Invalid, "帳號已存在", nil)
		return
	}
}

func Login(c *gin.Context) {
	var login LoginData
	// 使用者傳帳密進來，解析到struct
	err := c.ShouldBindJSON(&login)
	if err != nil {
		log.Printf("login BindJson Error : %s", err.Error())
		ErrMsg(c, Code_Param_Invalid, "帳號密碼格式錯誤", nil, err)
		return
	}
	// 確認帳號是否存在
	// Exists查詢key是否存在，回傳true 或是false
	check, err := redis.Client.Exists(login.Username).Result()
	if err != nil {
		log.Printf("Redis EXists Error :%s", err.Error())
		ErrMsg(c, Code_DB_Conn, "Redis EXists", nil, err)
		return
	}
	// check = 0，資料不存在，check = 1，資料已存在
	if check == 1 {
		value, err := redis.Client.Get(login.Username).Result()
		if err != nil {
			log.Printf("select Error :%s", err.Error())
			ErrMsg(c, Code_DB_Conn, "Redis EXists", nil, err)
			return
		}
		UserData = User{Username: login.Username, Password: value}
		// 資料庫取出的密碼 與 使用者輸入的密碼 驗證
		err = bcrypt.CompareHashAndPassword([]byte(UserData.Password), []byte(login.Password))
		if err != nil {
			Msg(c, Code_Param_Invalid, "密碼錯誤", nil)
			return
		}
		// 密碼對就生成Session、UserName加進Token回傳
		hash, err := bcrypt.GenerateFromPassword([]byte(UserData.Username), bcrypt.DefaultCost)
		if err != nil {
			ErrMsg(c, Code_Param_Invalid, "參數無效", nil, err)
			return
		}
		sessionID := string(hash)
		token, err := util.GenToken(sessionID, UserData.Username)
		if err != nil {
			ErrMsg(c, Code_Param_Invalid, "Token創建失敗", nil, err)
			return
		}
		// 把Token傳到前端cookie
		// cookie名字、cookie內容、存活時間，設定-1 = 刪除、path使用路經、作用host、是否只能https協議發送到服務端、HttpOnly=true 就不能被js獲取到、
		c.SetCookie("token", token, 600, "/", "localhost", false, true)
		// 創好的Token當Key，session ID當Value存進Redis
		err = redis.Client.Set(sessionID, UserData.Username, time.Minute*60).Err()
		if err != nil {
			log.Printf("Save Token Error :%s", err.Error())
			ErrMsg(c, Code_DB_Conn, "DB Error", nil, err)
			return
		}
		OkMsg(c, UserData.Username, token)
		log.Info().Msg(UserData.Username + "登錄成功")
	} else {
		Msg(c, Code_Param_Invalid, "查無此帳號", nil)
		return
	}
}

func Logout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		ErrMsg(c, Code_Param_Invalid, "請先登錄帳號", nil, err)
		return
	}
	c.SetCookie("token", token, -1, "/", "localhost", false, true)
	sessionid, ok := c.Get("sessionid")
	if !ok {
		Msg(c, Code_Param_Invalid, "驗證失敗", nil)
		return
	}
	err = redis.Client.Del(sessionid.(string), token).Err()
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "redis", nil, err)
		return
	}
	// err = redis.Client.Del(token).Err()
	// if err != nil {
	// 	Msg(c, Code_DB_Conn, "redis", nil, err)
	// 	return
	// }
	OkMsg(c, "登出成功", nil)
}

func LoginStruct(c *gin.Context) {
	var ok bool
	u, ok := c.Get("username")
	if !ok {
		Msg(c, Code_Param_Invalid, "", nil)
		return
	}
	username, ok := u.(string)
	if !ok {
		Msg(c, Code_Param_Invalid, "", nil)
		return
	}
	OkMsg(c, username, nil)
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
