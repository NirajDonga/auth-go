package httpserver

import (
	"go-auth/internal/app"
	"go-auth/internal/middleware"
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

	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthRequired(a.Config.JWTSecret))
	// Accessible to any authenticated user
	userRoutes.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, user!"})
	})

	// Admin-only route
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthRequired(a.Config.JWTSecret), middleware.RequireAdmin())
	adminRoutes.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, admin!"})
	})

	return r
}
