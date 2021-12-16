package util

import (
	"errors"
	"log"
	"strings"
	"time"
	"token/db"

	"github.com/golang-jwt/jwt"
)

type TokenData struct {
	Token string `json:"token"`
}
type MyClaims struct {
	SessionID string `json:"sessionid"`
	Username  string `json:"username"`
	jwt.StandardClaims
}

var TokenExpireDuration = time.Minute * 10
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

func GetUsername(username string) (bool, error) {
	var user db.User
	err := db.DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		log.Printf("select Error :%s", err.Error())
		user = db.User{}
		return false, err
	}
	return true, nil
}

// func bcryptPassword(data string){
// 	hash, err := bcrypt.GenerateFromPassword([]byte(UserData.Username), bcrypt.DefaultCost)
// 		if err != nil {
// 			ErrMsg(c, Code_Param_Invalid, "參數無效", nil, err)
// 			return
// 		}
// 		sessionID := string(hash)
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
// 		// ErrMsg(c, Code_DB_Conn, "Redis", nil, err)
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

// func GetUsername(c *gin.Context, username string) (interface{}, error) {
// 	var user controller.User
// 	// Exists查詢key是否存在，回傳true 或是false
// 	sel, err := util.RedisExists(c, "username")
// 	if err != nil {
// 		log.Printf("Error : %s", err.Error())
// 		return nil , nil
// 	}
// 	if sel == 1 {
// 		value, err := redis_DB.Client.Get(username).Result()
// 		if err != nil {
// 			log.Printf("select Error :%s", err.Error())
// 			return nil, err
// 		} else {
// 			user = controller.User{Username: username, Password: value}
// 			return user, nil
// 		}
// 	} else {
// 		log.Printf("查無此帳號")
// 	}
// 	return nil, nil
// }
