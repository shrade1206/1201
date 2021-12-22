package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"todoList/db"
	"todoList/util"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	Totle     int64 `json:"totle"`
	TotlePage int   `json:"totlepage"`
	// PageSize  int `json:"pagesize"`
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
		util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	}
	todos := util.FindAll(c)
	util.OkMsg(c, "新增成功", todos)

}

// 取得全部資料
func Get(c *gin.Context) {
	var todos []db.Todo
	err := db.DB.Find(&todos).Error
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	} else {
		util.OkMsg(c, "取得資料", todos)
		// c.JSON(http.StatusOK, todos)
	}
}

// 取得頁數
func GetPage(c *gin.Context) {
	DDB := db.DB
	var pages []db.Todo
	var p PageData
	var totle int64
	pageSize := 4
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		util.ErrMsg(c, util.Code_Param_Invalid, "參數無效", nil, err)
		return
	}
	err = DDB.Find(&pages).Count(&totle).Error
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	}
	strInt64 := strconv.FormatInt(totle, 10)
	t, err := strconv.Atoi(strInt64)
	if err != nil {
		return
	}

	totlePage := (t / pageSize)
	if t%pageSize == 0 {
		p = PageData{Totle: totle, TotlePage: totlePage}
	} else {
		a := totlePage + 1
		p = PageData{Totle: totle, TotlePage: a}
	}

	if page > 0 && pageSize > 0 {
		DDB = db.DB.Limit(pageSize).Offset((page - 1) * pageSize)
	}
	err = DDB.Find(&pages).Error
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Code": util.Code_ok,
		"Msg":  "",
		"Data": pages,
		"Page": p,
	})
}

// 更新資料
func Put(c *gin.Context) {
	var todo db.Todo
	id, ok := c.Params.Get("id")
	if !ok {
		util.Msg(c, util.Code_Param_Invalid, "id無效", nil)
		return
	}
	err := db.DB.Where("id =?", id).First(&todo).Error
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "資料更新錯誤", nil, err)
		return
	}
	err = c.BindJSON(&todo)
	if err != nil {
		util.ErrMsg(c, util.Code_Param_Invalid, "格式錯誤", nil, err)
		return
	}
	err = db.DB.Save(&todo).Error
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "資料保存錯誤", nil, err)
		return
	}
	todos := util.FindAll(c)
	util.OkMsg(c, "更新成功", todos)
}

// 刪除資料
func Del(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		util.Msg(c, util.Code_Param_Invalid, "參數無效", nil)
		return
	}
	err := db.DB.Where("id = ?", id).Delete(db.Todo{}).Error
	if err != nil {
		util.ErrMsg(c, util.Code_DB_Conn, "資料存取錯誤", nil, err)
		return
	}
	todos := util.FindAll(c)
	util.OkMsg(c, "刪除成功", todos)
}

func Router(c *gin.Context) {
	path := c.Request.URL.Path
	method := c.Request.Method
	fmt.Println(path)
	fmt.Println(method)
	//檢查path的開頭使是否為"/"
	if strings.HasPrefix(path, "/") {
		fmt.Println("ok")
	}
}
