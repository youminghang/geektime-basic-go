package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

// 专门用于 JWT 的代码

type UserClaims struct {
	// 我们只需要放一个 user id 就可以了
	Id        int64
	UserAgent string
	jwt.RegisteredClaims
}

// AccessTokenKey 因为 JWT Key 不太可能变，所以可以直接写成常量
// 也可以考虑做成依赖注入
var AccessTokenKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm")
var refreshTokenKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixA")

type jwtHandler struct {
}

func (h *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id:        uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 演示目的设置为一分钟过期
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			// 在压测的时候，要将过期时间设置更长一些
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	})
	tokenStr, err := token.SignedString(AccessTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h *jwtHandler) setRefreshToken(ctx *gin.Context,
	uid int64) error {
	rc := RefreshClaims{
		uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置为七天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	refreshTokenStr, err := refreshToken.SignedString(refreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

func ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}
	// SplitN 的意思是切割字符串，但是最多 N 段
	// 如果要是 N 为 0 或者负数，则是另外的含义，可以看它的文档
	authSegments := strings.SplitN(authCode, " ", 2)
	if len(authSegments) != 2 {
		// 格式不对
		return ""
	}
	return authSegments[1]
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	uid int64
	jwt.RegisteredClaims
}
