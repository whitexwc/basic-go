package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/whitexwc/basic-go/webook/internal/web"
	"log"
	"net/http"
	"strings"
	"time"
)

// LoginJWTMiddlewareBuilder JWT 登陆校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用go的方式编码解码
	gob.Register(time.Now())
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		// 现在使用 jwt 来校验
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登陆
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没登陆，有人瞎搞
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("BXRuAoqzeb4Tn6VjF1qcoUgntV0VEwq2"), nil
		})
		/*token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("BXRuAoqzeb4Tn6VjF1qcoUgntV0VEwq2"), nil
		})*/
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 刷新过期时间
		// 每10s钟刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			// 再生成一次token
			tokenStr, err = token.SignedString([]byte("BXRuAoqzeb4Tn6VjF1qcoUgntV0VEwq2"))
			if err != nil {
				//记录日志
				log.Println("jwt 续约失败", err)
			}
			c.Header("x-jwt-token", tokenStr)
		}

		// 将解析的结果塞进context
		c.Set("claims", claims)
		c.Set("userId", claims.Uid)
	}
}
