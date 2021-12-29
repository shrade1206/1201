package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"todoList/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type API_Error struct {
	Code int
	Msg  string
	Data interface{}
}

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

const Code_ok = 1
const Code_Param_Invalid = 2
const Code_DB_Conn = 3

// Token----------------------------------------
type MyClaims struct {
	SessionID string `json:"sessionid"`
	Username  string `json:"username"`
	jwt.StandardClaims
}

type TokenData struct {
	MySecret string
}

var Tk TokenData

// 解析Token ---------------------------------------------
func ParseToken(tokenString string) (*MyClaims, error) {
	file, err := os.Open("./config/TokenData.json")
	if err != nil {
		return nil, err
	}
	data := json.NewDecoder(file)
	err = data.Decode(&Tk)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Tk.MySecret), nil
	})
	if err != nil {
		expired := strings.Contains(err.Error(), "token is expired")
		if expired {
			return token.Claims.(*MyClaims), err
		}
		return nil, err
	}
	// 驗證claims正確就回傳
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid Token")
}

func FindAll(c *gin.Context, username string) interface{} {
	var todos []db.Todo
	err := db.DB.Where("user_id = ?", username).Find(&todos).Error
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "資料存取錯誤", nil, err)
		return nil
	}
	return todos
}

func GetUserName(c *gin.Context) string {
	u, ok := c.Get("username")
	if !ok {
		Msg(c, Code_Param_Invalid, "", nil)
		return ""
	}
	username := u.(string)
	return username
}
