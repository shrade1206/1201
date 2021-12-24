package util

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"token/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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
const Code_Session_Invalid = 4

// ---------------------------------------------

type MyClaims struct {
	SessionID string `json:"sessionid"`
	Username  string `json:"username"`
	jwt.StandardClaims
}

var TokenExpireDuration = time.Minute * 60
var MySecret = []byte("jiyuu")

// 生成token ---------------------------------------------
func GenToken(SessionID, username string) (string, error) {
	t := MyClaims{
		SessionID,
		username, // 自訂Header
		jwt.StandardClaims{ // 設定payload
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "Larry",
		},
	}
	// 選擇編碼模式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	// 用指定的SecretKey加密獲得Token字串
	return token.SignedString(MySecret)
}

// 解析Token ---------------------------------------------
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
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

// 檢查使用者名稱 ------------------------------------------
func CheckUsername(username string) (bool, error) {
	var user db.User
	err := db.DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return false, err
	}
	if user.Username == username {
		return true, nil
	}
	return false, nil
}

// 生成加密亂碼 --------------------------------------------
func BcryptPassword(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	bcryptString := string(hash)
	return bcryptString, nil
}

// 取得使用者資料 -------------------------------------------
func GetUserData(username string) (interface{}, error) {
	var user db.User
	err := db.DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 讀取DsnConfig ------------------------------------------
// type RedisDsnData struct {
// 	Addr        string
// 	Password    string
// 	DB          int
// 	PoolSize    int
// 	PoolTimeout int
// 	MaxConnAge  int
// }

// func RedisDsn(config string) (interface{}, error) {
// 	var dsn RedisDsnData
// 	file, err := os.Open(config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	data := json.NewDecoder(file)
// 	err = data.Decode(&dsn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dsn, nil
// }

// 讀取DsnConfig ---------------------------------------------
// func DsnGet() {
// 	file, err := os.Open("./config/redis.json")
// 	if err != nil {
// 		return
// 	}
// 	var dsn redis.Dsn
// 	data := json.NewDecoder(file)
// 	err = data.Decode(&dsn)
// 	if err != nil {
// 		return
// 	}
// }

// Redis驗證Key是否存在 ---------------------------------------------
// func RedisExists(c *gin.Context, username string) (int64, error) {
// 	check, err := redis.Client.Exists(username).Result()
// 	if err != nil {
// 		return check, err
// 	}
// 	return check, nil
// }

// Redis設定Key.Value ---------------------------------------------
// func RedisSet(c *gin.Context, key string, value interface{}, expiration time.Duration) (string, error) {
// 	SaveData, err := redis.Client.Set(key, value, expiration).Result()
// 	if err != nil {
// 		ErrMsg(c, Code_DB_Conn, "Redis", nil, err)
// 		return SaveData, err
// 	}
// 	return SaveData, nil
// }

// func Test(c *gin.Context) {
// 	a, _ := GetUsername(c, "jiyuusama")
// 	b := a.(db.User)
// 	fmt.Println(b.Username, b.ID, b.Password)
// }
