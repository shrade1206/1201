package util

import (
	"net/http"
	"todoList/db"

	"github.com/gin-gonic/gin"
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

func OkMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, API_Error{
		Code: Code_ok,
		Msg:  msg,
		Data: data,
	})
}

const Code_ok = 1
const Code_Param_Invalid = 2
const Code_DB_Conn = 3

var todos []db.Todo

func FindAll(c *gin.Context) interface{} {
	err := db.DB.Find(&todos).Error
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "資料存取錯誤", nil, err)
		return nil
	}
	return todos
}
