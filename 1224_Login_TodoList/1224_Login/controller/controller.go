package controller

import (
	"fmt"
	"time"

	"strings"
	"token/db"
	"token/redis"
	"token/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterData struct {
	Username string `json:"username" form:"username" binding:"required,min=6,max=12"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
}

type LoginData struct {
	Username string `json:"username" form:"username" binding:"required,min=6,max=12"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=12"`
}

//註冊
func Register(c *gin.Context) {
	var user db.User
	var reg RegisterData
	// 使用者傳帳密進來，解析到struct
	err := c.ShouldBindJSON(&reg)
	if err != nil {
		util.Msg(c, util.Code_Param_Invalid, "帳號密碼格式錯誤", nil)
		return
	}
	// 確認帳號是否存在，不存在就Create
	checkOK, err := util.CheckUsername(reg.Username)
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤#1", nil, err)
		return
	}
	if !checkOK {
		// 密碼加密
		regPassword, err := util.BcryptPassword(reg.Password)
		if err != nil {
			util.ErrMsg(c, util.Code_Param_Invalid, "資料錯誤", nil, err)
			return
		}
		// 創建帳密key、value
		user = db.User{Username: reg.Username, Password: regPassword}
		err = db.DB.Create(&user).Error
		if err != nil {
			util.ErrMsg(c, util.Code_DB_Conn, "註冊失敗，請重新嘗試", nil, err)
			return
		}
		util.OkMsg(c, "註冊成功", nil)
	} else {
		// 帳號已存在
		util.Msg(c, util.Code_Param_Invalid, "帳號已存在", nil)
		return
	}
}

// 登入
func Login(c *gin.Context) {
	var login LoginData
	// 使用者傳帳密進來，解析到struct
	err := c.ShouldBindJSON(&login)
	if err != nil {
		util.Msg(c, util.Code_Param_Invalid, "帳號密碼格式錯誤", nil)
		return
	}
	// 確認帳號是否存在
	checkOK, err := util.CheckUsername(login.Username)
	if err != nil {
		util.ErrMsg(c, util.Code_Param_Invalid, "資料存取錯誤#2", nil, err)
		return
	}
	// check = 0，資料不存在，check = 1，資料已存在
	if checkOK {
		data, err := util.GetUserData(login.Username)
		if err != nil {
			util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤#1", nil, err)
			return
		}
		user := data.(db.User)
		// 資料庫取出的密碼 與 使用者輸入的密碼 驗證
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
		if err != nil {
			util.Msg(c, util.Code_Param_Invalid, "密碼錯誤", nil)
			return
		}
		// 密碼對就生成Session、UserName加進Token回傳
		sessionID, err := util.BcryptPassword(user.Username)
		if err != nil {
			util.ErrMsg(c, util.Code_Param_Invalid, "參數無效", nil, err)
			return
		}
		token, err := util.GenToken(sessionID, user.Username)
		if err != nil {
			util.ErrMsg(c, util.Code_Param_Invalid, "Token創建失敗", nil, err)
			return
		}
		// 把Token傳到前端cookie
		// cookie名字、cookie內容、存活時間，設定-1 = 刪除、path使用路經、作用host、是否只能https協議發送到服務端、HttpOnly=true 就不能被js獲取到、
		c.SetCookie("token", token, 3600, "/", "localhost", false, true)
		// 創好的Token當Key，session ID當Value存進Redis
		err = redis.Client.Set(sessionID, user.Username, time.Minute*60).Err()
		if err != nil {
			util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤#2", nil, err)
			return
		}
		util.OkMsg(c, user.Username, token)
	} else {
		util.Msg(c, util.Code_Param_Invalid, "查無此帳號", nil)
		return
	}
}

// 登出
func Logout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		util.ErrMsg(c, util.Code_Param_Invalid, "請先登錄帳號", nil, err)
		return
	}
	c.SetCookie("token", token, -1, "/", "localhost", false, true)
	sessionid, ok := c.Get("sessionid")
	if !ok {
		util.Msg(c, util.Code_Param_Invalid, "驗證失敗", nil)
		return
	}
	err = redis.Client.Del(sessionid.(string), token).Err()
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "redis", nil, err)
		return
	}
	util.OkMsg(c, "登出成功", nil)
}

// 確認當前使用者名稱
func LoginStruct(c *gin.Context) {
	var ok bool
	u, ok := c.Get("username")
	if !ok {
		util.Msg(c, util.Code_Param_Invalid, "", nil)
		return
	}
	username, ok := u.(string)
	if !ok {
		util.Msg(c, util.Code_Param_Invalid, "", nil)
		return
	}
	util.OkMsg(c, username, nil)
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
