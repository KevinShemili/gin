package authentication

import (
	"gin/api/middleware"
	"gin/application/repository/contracts"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(r *gin.Engine, controller *AuthenticationController, UnitOfWork contracts.IUnitOfWork) {

	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.POST("/refresh", controller.Refresh)
		auth.POST("/logout", middleware.AuthenticationMiddleware(UnitOfWork), controller.LogOut)
	}
}
