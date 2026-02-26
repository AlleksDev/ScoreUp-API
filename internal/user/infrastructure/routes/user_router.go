package routes

import (
	//"github.com/AlleksDev/ScoreUp-API/internal/middleware"
	"github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, createUserCtrl *controllers.CreateUserController,
	jwtSecret string) {
	g := r.Group("users")
	{
		g.POST("", createUserCtrl.Handle)
	}
	/* gPrivate := r.Group("users")
	gPrivate.Use(middleware.AuthMiddleware(jwtSecret))
	{
		gPrivate.GET("/get/:id", getUserCtrl.Handle)
		gPrivate.GET("/username/:username", getByUsernameCtrl.Handle)
		gPrivate.GET("/search", searchUsersCtrl.Handle)
		gPrivate.GET("/profile", getProfileCtrl.Handle)
		gPrivate.GET("/invitations", getPendingInvitationsCtrl.Handle)
		gPrivate.PATCH("/update", updateUserCtrl.Handle)
		gPrivate.DELETE("/delete", deleteUserCtrl.Handle)
	} */
}