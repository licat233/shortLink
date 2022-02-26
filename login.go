package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/licat233/goutil/readfile"
)

type Admin struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

//自定义一个字符串
var jwtkey = []byte("planttitle.com")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

type TokenInfo struct {
	Token     string    `json:"Token"`
	ExpiresAt time.Time `json:"ExpiresAt"`
}

type LoginReq struct {
	UserName  string `json:"Username"`
	Password  string `json:"Password"`
	AutoLogin bool   `json:"AutoLogin"`
}

func (a *App) loginVerify(ctx *gin.Context) {
	req := &LoginReq{}
	ctx.BindJSON(req)
	//為了方便隨時更改賬號密碼，而不用重啟服務
	readfile.YamlConfig("./config.yaml", &a.Config, func(err error) {
		if err != nil {
			panic(err)
		}
	})
	if req.UserName != a.Config.Admin.Username || req.Password != a.Config.Admin.Password {
		ctx.JSON(200, gin.H{
			"code":    401,
			"message": "賬號/密碼錯誤",
		})
		return
	}
	token, err := setting(ctx, req.AutoLogin)
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": "jwt服务异常",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "登錄成功",
		"data":    token,
	})
}

func (a *App) Loginpage(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "登錄頁面",
	})
}

//颁发token
func setting(ctx *gin.Context, AutoLogin bool) (*TokenInfo, error) {
	expireTime := time.Now().Add(12 * time.Hour)
	if AutoLogin {
		expireTime = time.Now().Add(7 * 24 * time.Hour)
	}
	claims := &Claims{
		UserId: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwttoken.SignedString(jwtkey)
	if err != nil {
		return nil, err
	}
	return &TokenInfo{Token: token, ExpiresAt: expireTime}, err
}

//解析token
func getting(ctx *gin.Context) *Claims {
	tokenString := ctx.GetHeader("Authorization")
	// fmt.Println(tokenString)
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "權限不足,請先登錄"})
		ctx.Abort()
		return nil
	}

	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "權限不足,請重新登錄"})
		ctx.Abort()
		return nil
	}
	return claims
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}
