package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"web_socket/src/online_medicine_store/define"
)

func GenerateToken(claim define.JwtClaims) (string, error) {
	// Create the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ss, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func AnalysisToken(tokens string) (*define.JwtClaims, error) {
	// Create the claims
	token, err := jwt.ParseWithClaims(tokens, &define.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*define.JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func JsonRecord(by, to uint32, msg string) string {
	s := fmt.Sprint("{\"by\":", by, ",\"to\":", to, ",\"msg\":\"", msg, "\"}")
	return s
}
