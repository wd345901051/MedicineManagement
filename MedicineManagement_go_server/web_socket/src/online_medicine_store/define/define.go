package define

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	Root uint32 `json:"root"` // 0,普通用户,1,房主
	jwt.RegisteredClaims
}

type Action2 struct {
	Id  uint32 `json:"id"`
	Msg string `json:"msg"`
}

/*
{
	id:2,data:"{id:userid,msg:"hello"}"
}
*/
