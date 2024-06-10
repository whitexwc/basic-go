package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range l.paths {
			if c.Request.URL.Path == path {
				return
			}
		}
		// 不需要登陆校验的
		// if c.Request.URL.Path == "/users/login" || c.Request.URL.Path == "/users/signup" {
		// 	return
		// }
		sess := sessions.Default(c)
		id := sess.Get("userId")
		if id == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
