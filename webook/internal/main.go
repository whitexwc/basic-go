package main

import (
	"github.com/gin-contrib/sessions/redis"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/whitexwc/basic-go/webook/internal/repository"
	"github.com/whitexwc/basic-go/webook/internal/repository/dao"
	"github.com/whitexwc/basic-go/webook/internal/service"
	"github.com/whitexwc/basic-go/webook/internal/web"
	"github.com/whitexwc/basic-go/webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := initDB()
	server := initServer()

	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"authorization", "content-type"},
		//ExposeHeaders:    []string{},
		// 是否允许带cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// 允许所有本地开发环境
			if strings.Contains(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 步骤1
	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("BXRuAoqzeb4Tn6VjF1qcoUgntV0VEwq2"),
	//	[]byte("7BS1f8ZqOaPuo7IBo3gJtOQzhh2P3NMX"))
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("BXRuAoqzeb4Tn6VjF1qcoUgntV0VEwq2"), []byte("7BS1f8ZqOaPuo7IBo3gJtOQzhh2P3NMX"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))

	// 步骤3
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").IgnorePaths("/users/login").
		IgnorePaths("/users/profile").IgnorePaths("/users/edit").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	du := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(du)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		// 只在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
