package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	// 用go的方式编码解码
	gob.Register(time.Now())
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

		// 如何知道一分钟已经过去了
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 30,
		})
		now := time.Now()
		// 说明是第一次登陆,还没有刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		// 如果有update_time
		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			sess.Save()
		}

	}
}
