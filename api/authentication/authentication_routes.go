package authentication

import (
	"gin/api/middleware"
	"gin/application/repository/contracts"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(r *gin.Engine, controller *AuthenticationController, UnitOfWork contracts.IUnitOfWork) {
	jwtSecretKey := os.Getenv("SECRET_JWT")

	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.POST("/logout", middleware.AuthenticationMiddleware(UnitOfWork, jwtSecretKey), controller.LogOut)
	}
}
