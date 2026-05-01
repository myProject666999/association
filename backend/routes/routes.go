package routes

import (
	"association/controllers"
	"association/middleware"
	"association/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/user/me", controllers.GetCurrentUser)

			users := auth.Group("/users")
			users.Use(middleware.RoleMiddleware(string(models.RoleUniversityAdmin), string(models.RoleDeptAdmin)))
			{
				users.GET("", controllers.GetUsers)
				users.GET("/:id", controllers.GetUserByID)
				users.POST("", controllers.CreateUser)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
				users.PUT("/:id/role", controllers.UpdateUserRole)
				users.PUT("/:id/status", controllers.ToggleUserStatus)
			}

			clubs := auth.Group("/clubs")
			{
				clubs.GET("", controllers.GetClubs)
				clubs.GET("/:id", controllers.GetClubByID)
				clubs.POST("", controllers.CreateClub)
				clubs.PUT("/:id", controllers.UpdateClub)
				clubs.DELETE("/:id", controllers.DeleteClub)
				clubs.POST("/:id/apply", controllers.ApplyJoinClub)
				clubs.GET("/:id/members", controllers.GetClubMembers)
				clubs.PUT("/:id/members/:member_id/review", controllers.ReviewMemberApplication)
				clubs.DELETE("/:id/members/:member_id", controllers.RemoveMember)
			}

			activities := auth.Group("/activities")
			{
				activities.GET("", controllers.GetActivities)
				activities.GET("/:id", controllers.GetActivityByID)
				activities.POST("", controllers.CreateActivity)
				activities.PUT("/:id", controllers.UpdateActivity)
				activities.DELETE("/:id", controllers.DeleteActivity)
				activities.PUT("/:id/review", controllers.ReviewActivity)
				activities.POST("/:id/register", controllers.RegisterActivity)
				activities.GET("/:id/registrations", controllers.GetActivityRegistrations)
				activities.PUT("/:id/registrations/:registration_id/review", controllers.ReviewRegistration)
				activities.POST("/:id/comments", controllers.AddComment)
				activities.GET("/:id/comments", controllers.GetActivityComments)
				activities.DELETE("/:id/comments/:comment_id", controllers.DeleteComment)
			}

			profile := auth.Group("/profile")
			{
				profile.GET("", controllers.GetMyProfile)
				profile.PUT("", controllers.UpdateMyProfile)
				profile.PUT("/password", controllers.ChangeMyPassword)
				profile.GET("/clubs", controllers.GetMyClubs)
				profile.GET("/activities", controllers.GetMyActivities)
				profile.GET("/organized-activities", controllers.GetMyOrganizedActivities)
				profile.GET("/club-applications", controllers.GetMyClubApplications)
				profile.GET("/activity-applications", controllers.GetMyActivityApplications)
			}
		}
	}

	return r
}
