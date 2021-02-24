package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Ext struct {
	Is_admin bool   `json:"is_admin"`
	Conn_id  string `json:"conn_id"`
}

type Claims struct {
	Kid      string `json:"kid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
	Ext      Ext    `json:"ext"`
	jwt.StandardClaims
}

var jwtSecret = []byte("ok")

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		"0",
		username,
		password,
		true,
		Ext{},
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func ParseTokenUnverified(token string) (*Claims, error) {
	tokenClaims, _, err := new(jwt.Parser).ParseUnverified(token, &Claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok {
		claims.Kid = tokenClaims.Header["kid"].(string)
		return claims, nil
	} else {
		return nil, nil
	}
}
