// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/adityjoshi/avinyaa/consumer"
// 	"github.com/adityjoshi/avinyaa/controllers"
// 	"github.com/adityjoshi/avinyaa/database"
// 	"github.com/adityjoshi/avinyaa/initiliazers"
// 	kafkamanager "github.com/adityjoshi/avinyaa/kafka/kafkaManager"
// 	"github.com/adityjoshi/avinyaa/routes"
// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-contrib/sessions"
// 	"github.com/gin-contrib/sessions/cookie"
// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// )

// var km *kafkamanager.KafkaManager

// func init() {
// 	initiliazers.LoadEnvVariable()
// }

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Initialize database
// 	database.InitDatabase()
// 	defer database.CloseDatabase()
// 	database.InitializeRedisClient()

// 	// Kafka initialization
// 	northBrokers := []string{"localhost:9092"}
// 	southBrokers := []string{"localhost:9092"}

// 	var err error
// 	km, err = kafkamanager.NewKafkaManager(northBrokers, southBrokers) // Initialize Kafka Manager
// 	if err != nil {
// 		log.Fatal("Failed to initialize Kafka Manager:", err)
// 	}

// 	region := "north" // Dynamically set this based on your application’s logic
// 	consumer.StartConsumer(region)
// 	select {}

// 	// brokers := []string{"localhost:9092"} // Add more brokers if needed

// 	// // Start Kafka consumers for North and South regions
// 	// go kafkaconsumer.StartConsumer(brokers, "north") // Consumer for North region
// 	// go kafkaconsumer.StartConsumer(brokers, "south") // Consumer for South region
// 	// brokers := []string{"localhost:9092"}
// 	// region := "north" // or "south"
// 	// kafkaconsumer.StartConsumer(brokers, region)

// 	router := gin.Default()
// 	go controllers.SubscribeToPaymentUpdates()
// 	go controllers.SubscribeToHospitalizationUpdates()

// 	// CORS setup
// 	router.Use(setupCORS())

// 	// Session setup
// 	store := cookie.NewStore([]byte("secret"))
// 	router.Use(sessions.Sessions("session", store))

// 	// Route definitions
// 	routes.UserRoutes(router)
// 	routes.UserInfoRoutes(router)
// 	routes.HospitalAdmin(router, km)

// 	// Start server
// 	server := &http.Server{
// 		Addr:    ":2426",
// 		Handler: router,
// 	}

// 	log.Println("Server is running at :2426...")
// 	server.ListenAndServe()
// }

// func setupCORS() gin.HandlerFunc {
// 	config := cors.DefaultConfig()
// 	config.AllowOrigins = []string{"http://localhost:3000"}
// 	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
// 	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
// 	config.AllowCredentials = true

// 	return cors.New(config)
// }

package main

import (
	"log"
	"net/http"

	"github.com/adityjoshi/avinyaa/consumer"
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
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	database.InitDatabase()
	defer database.CloseDatabase()
	database.InitializeRedisClient()

	// Initialize Kafka Manager with regional brokers
	northBrokers := []string{"localhost:9092"}
	southBrokers := []string{"localhost:9092"}
	var err error
	km, err = kafkamanager.NewKafkaManager(northBrokers, southBrokers)
	if err != nil {
		log.Fatal("Failed to initialize Kafka Manager:", err)
	}

	// Start consumers for each region in separate goroutines
	regions := []string{"north", "south"}
	for _, region := range regions {
		go func(r string) {
			log.Printf("Starting Kafka consumer for region: %s\n", r)
			consumer.StartConsumer(r) // Start each region's consumer
		}(region)
	}

	// Start consumers for each region in separate goroutines

	// Start additional goroutines for other tasks
	go controllers.SubscribeToPaymentUpdates()
	go controllers.SubscribeToHospitalizationUpdates()
	go controllers.SubscribeToHospitaliztionUpdates()

	// Setup the HTTP server with Gin
	router := gin.Default()
	router.Use(setupCORS())
	setupSessions(router)
	setupRoutes(router)

	// Start server
	server := &http.Server{
		Addr:    ":2426",
		Handler: router,
	}
	log.Println("Server is running at :2426...")

	// Keep main function running indefinitely
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	select {} // Prevent main from exiting
}

// setupCORS configures CORS settings
func setupCORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowCredentials = true
	return cors.New(config)
}

// setupSessions configures session management
func setupSessions(router *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))
}

// setupRoutes defines all application routes
func setupRoutes(router *gin.Engine) {
	routes.UserRoutes(router)
	routes.UserInfoRoutes(router)
	routes.HospitalAdmin(router, km)
}
