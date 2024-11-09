package routes

import (
	"github.com/adityjoshi/avinyaa/controllers"
	"github.com/adityjoshi/avinyaa/middleware"
	"github.com/gin-gonic/gin"
)

func UserInfoRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/updatePatientInfo/:id", middleware.AuthRequired("Patient", ""), middleware.OtpAuthRequireed, controllers.AddPatientDetails)
	// incomingRoutes.GET("/getPatientId/:id", middleware.AuthUser(), controllers.GetPatientDetails)

}
