package controller

import (
	"log"
	"net/http"
	"strconv"
	"todoList/db"

	"github.com/gin-gonic/gin"
)

type API_Error struct {
	Code int
	Msg  string
	Data interface{}
}

const Code_ok = 1
const Code_Param_Invalid = 2
const Code_DB_Conn = 3

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

// 新增一筆資料
func Insert(c *gin.Context) {
	var todo db.Todo
	err := c.BindJSON(&todo)
	if err != nil {
		log.Printf("BindJson Error : %s", err.Error())
		return
	}
	err = db.DB.Create(&todo).Error
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	} else {
		OkMsg(c, todo.Title, nil)
	}
}

// 取得全部資料
func Get(c *gin.Context) {
	var todos []db.Todo
	err := db.DB.Find(&todos).Error
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	} else {
		OkMsg(c, "取得資料", todos)
		// c.JSON(http.StatusOK, todos)
	}
}

// 取得頁數
func GetPage(c *gin.Context) {
	DDB := db.DB
	var pages []db.Todo
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		ErrMsg(c, Code_Param_Invalid, "參數無效", nil, err)

		return
	}
	pageSize := 4
	if page > 0 && pageSize > 0 {
		DDB = db.DB.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	err = DDB.Find(&pages).Error
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	} else {
		OkMsg(c, "取得頁數", page)
		// c.JSON(http.StatusOK, pages)
	}
}

// 更新資料
func Put(c *gin.Context) {
	var todo db.Todo
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"Put Error": "ID無效",
		})
		return
	}
	err := db.DB.Where("id =?", id).First(&todo).Error
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "資料更新錯誤", nil, err)
		return
	}
	err = c.BindJSON(&todo)
	if err != nil {
		ErrMsg(c, Code_Param_Invalid, "格式錯誤", nil, err)
		return
	}
	err = db.DB.Save(&todo).Error
	if err != nil {
		ErrMsg(c, Code_DB_Conn, "資料保存錯誤", nil, err)
		return
	} else {
		OkMsg(c, "更新成功", todo)
		// c.JSON(http.StatusOK, gin.H{
		// 	"id":     todo.Id,
		// 	"title":  todo.Title,
		// 	"status": todo.Status,
		// })
	}
}

// 刪除資料
func Del(c *gin.Context) {
	id, ok := c.Params.Get("id")
	log.Printf("id: %v\n", id)
	if !ok {
		Msg(c, Code_Param_Invalid, "參數無效", nil)
		return
	}
	err := db.DB.Where("id = ?", id).Delete(db.Todo{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
		return
	} else {
		OkMsg(c, "刪除成功", nil)
	}
}
