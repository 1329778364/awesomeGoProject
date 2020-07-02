package util

import (
	"time"
	"userSystem/pkg/setting"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
}

func GenerateToken(userId string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(8 * time.Hour)
	claims := Claims{
		jwt.StandardClaims{
			Audience:  EncodeMD5(userId),  // 受众
			ExpiresAt: expireTime.Unix(),  // 失效时间
			Id:        EncodeMD5(userId),  // 编号
			IssuedAt:  time.Now().Unix(),  // 签发时间
			Issuer:    "awesomeGoProject", // 签发人
			NotBefore: time.Now().Unix(),  // 生效时间
			Subject:   "login",            // 主题
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(setting.AppSetting.JwtSecret))
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.AppSetting.JwtSecret), nil
	})
	if jwtToken != nil {
		if claims, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
			return claims, nil
		}
	}
	return nil, err
}
