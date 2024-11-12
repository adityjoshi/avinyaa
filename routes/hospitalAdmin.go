package routes

import (
	"time"

	"github.com/adityjoshi/avinyaa/controllers"
	kafkamanager "github.com/adityjoshi/avinyaa/kafka/kafkaManager"
	"github.com/adityjoshi/avinyaa/middleware"
	"github.com/gin-gonic/gin"
)

func HospitalAdmin(incomingRoutes *gin.Engine, km *kafkamanager.KafkaManager) {
	incomingRoutes.POST("/hospitaladmin", func(c *gin.Context) {
		c.Set("km", km)                      // Set km into the context
		controllers.RegisterHospitalAdmin(c) // Call the controller function
	})
	incomingRoutes.POST("/logindoc", controllers.DoctorLogin)
	incomingRoutes.POST("/adminLogin", middleware.RateLimiterMiddleware(2, time.Minute), controllers.AdminLogin)
	incomingRoutes.POST("/adminOtp", middleware.AuthRequired("Admin", ""), controllers.VerifyAdminOTP)
	//incomingRoutes.POST("/stafflogin", controllers.StaffLogin)
	incomingRoutes.POST("/staffotp", middleware.AuthRequired("Staff", ""), controllers.VerifyStaffOTP)

	//incomingRoutes.POST("/admit", middleware.AuthRequired("Staff", "Compounder"), controllers.AdmitPatientForHospitalization)
	incomingRoutes.POST("/compounder", controllers.CompounderLogin)
	incomingRoutes.POST("/markCompounder", middleware.AuthRequired("Staff", "Compounder"), controllers.MarkPatientAsHospitalized)
	incomingRoutes.GET("/get", middleware.AuthRequired("Staff", "Compounder"), controllers.GetRoomAssignments)
	incomingRoutes.POST("/registerhospital", middleware.AuthRequired("Admin", ""), func(c *gin.Context) {
		// Set km into the context (KafkaManager)
		c.Set("km", km)

		// Call the controller function
		controllers.RegisterHospital(c)
	})
	// incomingRoutes.GET("/gethospital/:id", controllers.GetHospital)
	incomingRoutes.POST("/doctor", middleware.AuthRequired("Admin", ""), controllers.RegisterDoctor)
	// incomingRoutes.GET("/getdoctor/:id", controllers.GetDoctor)
	// incomingRoutes.POST("/bookAppointment", controllers.CreateAppointment)
	incomingRoutes.POST("/markAppointment/:appointment_id", controllers.MarkAppointmentAsDone)

	adminRoutes := incomingRoutes.Group("/admin")
	adminRoutes.Use(middleware.AuthRequired("Admin", ""))
	{
		//adminRoutes.POST("/registerhospital", middleware.OtpAuthRequireed, controllers.RegisterHospital)
		//adminRoutes.GET("/gethospital/:id", middleware.OtpAuthRequireed, controllers.GetHospital)
		adminRoutes.POST("/doctor", middleware.OtpAuthRequireed, controllers.RegisterDoctor)
		adminRoutes.GET("/getdoctor/:id", middleware.OtpAuthRequireed, controllers.GetDoctor)
		adminRoutes.POST("/v", middleware.OtpAuthRequireed, func(c *gin.Context) {
			c.Set("km", km)
			controllers.CreateAppointment(c)
		})
		adminRoutes.POST("/registerStaff", middleware.OtpAuthRequireed, func(c *gin.Context) {
			c.Set("km", km)
			controllers.RegisterStaff(c)
		})
		adminRoutes.POST("/registerBeds", middleware.OtpAuthRequireed, controllers.AddBedType)
		adminRoutes.POST("/updateBeds", middleware.OtpAuthRequireed, controllers.UpdateTotalBeds)
		//adminRoutes.GET("/getBeds", middleware.OtpAuthRequireed, controllers.GetTotalBeds)
	}
	receptionistRoute := incomingRoutes.Group("/reception")
	receptionistRoute.Use(middleware.AuthRequired("Staff", "Compounder"))
	{
		incomingRoutes.POST("/stafflogin", controllers.StaffLogin)
		receptionistRoute.POST("/staffotp", controllers.VerifyStaffOTP)
		receptionistRoute.POST("/patientReg", func(c *gin.Context) {
			c.Set("km", km)
			controllers.RegisterPatient(c)
		})
		// receptionistRoute.POST("/admit", func(c *gin.Context) {
		// 	c.Set("km", km)
		// 	controllers.AdmitPatientForHospitalization(c)
		// })
	}
	compounderRoute := incomingRoutes.Group("/compounder")
	compounderRoute.Use(middleware.AuthRequired("Staff", "Compounder"))
	{
		incomingRoutes.POST("/compounderstafflogin", controllers.CompounderLogin)
		compounderRoute.POST("/admit", func(c *gin.Context) {
			c.Set("km", km)
			controllers.AdmitPatientForHospitalization(c)
		})
		compounderRoute.POST("/hospitalise", func(c *gin.Context) {
			c.Set("km", km)
			controllers.AdmitPatient(c)
		})
	}
}
