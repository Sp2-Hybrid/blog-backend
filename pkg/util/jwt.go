package util

import (
	jwt "github.com/dgrijalva/jwt-go"
	"go-gin-example/pkg/setting"
	"time"
)
var jwtSecret = []byte(setting.JwtSecret)

// 创建一个我们自己的声明
type Claims struct{
	// 自定义字段
	Username 	string 		`json:"username"`
	Password 	string 		`json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string)(string, error)  {
	nowTime := time.Now()
	// 到期时间
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		Username:       username,
		Password:       password,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt:expireTime.Unix(),
			// 签发人
			Issuer:"gin-blog",
		},
	}
	// 使用指定的签名方式创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string)(*Claims, error){
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil{
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid{
			return claims, nil
		}
	}
	return nil, err
}
