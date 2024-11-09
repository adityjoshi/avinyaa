package main

import (
	"log"
	"net/http"

	"github.com/adityjoshi/avinyaa/controllers"
	"github.com/adityjoshi/avinyaa/database"
	"github.com/adityjoshi/avinyaa/initiliazers"
	kafkamanager "github.com/adityjoshi/avinyaa/kafka/kafkaManager"
	"github.com/adityjoshi/avinyaa/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var km *kafkamanager.KafkaManager

func init() {
	initiliazers.LoadEnvVariable()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	database.InitDatabase()
	defer database.CloseDatabase()
	database.InitializeRedisClient()

	// Kafka initialization
	northBrokers := []string{"localhost:9092"}
	southBrokers := []string{"localhost:9092"}

	var err error
	km, err = kafkamanager.NewKafkaManager(northBrokers, southBrokers) // Initialize Kafka Manager
	if err != nil {
		log.Fatal("Failed to initialize Kafka Manager:", err)
	}

	router := gin.Default()
	go controllers.SubscribeToPaymentUpdates()
	go controllers.SubscribeToHospitalizationUpdates()

	// CORS setup
	router.Use(setupCORS())

	// Session setup
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	// Route definitions
	routes.UserRoutes(router)
	routes.UserInfoRoutes(router)
	routes.HospitalAdmin(router, km)

	// Start server
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
