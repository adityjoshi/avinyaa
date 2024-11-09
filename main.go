package main

import (
	"log"
	"net/http"

	"github.com/adityjoshi/avinya/controllers"
	"github.com/adityjoshi/avinyaa/database"
	"github.com/adityjoshi/avinyaa/initiliazers"
	"github.com/adityjoshi/avinyaa/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	initiliazers.LoadEnvVariable()
}
func main() {
	if err := godotenv.Load(); err != nil {

		log.Fatal("Error loading .env file")
	}
	region := "null"
	database.InitDatabase(region)
	defer database.CloseDatabase()
	database.InitializeRedisClient()

	router := gin.Default()
	go controllers.SubscribeToPaymentUpdates()
	go controllers.SubscribeToHospitalizationUpdates()
	// router.Use(cors.Default())
	router.Use(setupCORS())
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	routes.UserRoutes(router)
	routes.UserInfoRoutes(router)
	routes.HospitalAdmin(router)

	server := &http.Server{
		Addr:    ":2426",
		Handler: router,
	}

	log.Println("Server is running at :2426...")
	server.ListenAndServe()
}

func setupCORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowCredentials = true

	return cors.New(config)
}

// and folders will be like controllers, database, initilizares, middleware , routes, utils
