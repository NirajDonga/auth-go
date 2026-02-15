package httpserver

import (
	"go-auth/internal/app"
	"go-auth/internal/user"

	"github.com/gin-gonic/gin"
)

func NewRouter(a *app.App) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", health)

	userRepo := user.NewRepo(a.DB)
	userSvc := user.NewService(userRepo, a.Config.JWTSecret)
	userHandler := user.NewHandler(userSvc)

	// unauth routes -> public routes
	r.POST("/register", userHandler.RegisterHandler)
	r.POST("/login", userHandler.LoginHandler)

	return r
}
