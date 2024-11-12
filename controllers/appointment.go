package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adityjoshi/avinyaa/database"
	kafkamanager "github.com/adityjoshi/avinyaa/kafka/kafkaManager"
	"github.com/gin-gonic/gin"
)

// // CreateAppointment handles POST requests to create a new appointment
// func CreateAppointment(c *gin.Context) {
// 	var appointmentData struct {
// 		PatientID       uint      `json:"patient_id"`
// 		DoctorID        uint      `json:"doctor_id"`
// 		AppointmentDate time.Time `json:"appointment_date"`
// 		AppointmentTime time.Time `json:"appointment_time"`
// 		Description     string    `json:"description"`
// 	}

// 	if err := c.BindJSON(&appointmentData); err != nil {
// 		fmt.Println("Error binding JSON:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
// 		return
// 	}

// 	appointment := database.Appointment{
// 		PatientID:       appointmentData.PatientID,
// 		DoctorID:        appointmentData.DoctorID,
// 		AppointmentDate: appointmentData.AppointmentDate,
// 		AppointmentTime: appointmentData.AppointmentTime,
// 		Description:     appointmentData.Description,
// 	}

// 	if err := database.DB.Create(&appointment).Error; err != nil {
// 		fmt.Println("Error creating appointment:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
// 		return
// 	}
// 	var doctor database.Doctors
// 	if err := database.DB.Where("doctor_id = ?", appointmentData.DoctorID).First(&doctor).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctor information"})
// 		return
// 	}

// 	// Fetch patient details
// 	var patient database.Users
// 	if err := database.DB.Where("patient_id = ?", appointmentData.PatientID).First(&patient).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve patient information"})
// 		return
// 	}

// 	// Send appointment email to the patient
// 	bookingTime := time.Now().Format("2006-01-02 15:04:05") // Current time as booking time
// 	err := utils.SendAppointmentEmail(patient.Email, doctor.FullName, appointmentData.AppointmentDate.Format("2006-01-02"), appointmentData.AppointmentTime.Format("15:04"), bookingTime)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send appointment email"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Appointment created successfully", "appointment_id": appointment.AppointmentID})
// }

func CreateAppointment(c *gin.Context) {
	km, exists := c.Get("km")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "KafkaManager not found"})
		return
	}

	kafkaManager, ok := km.(*kafkamanager.KafkaManager)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid KafkaManager"})
		return
	}

	var appointmentData database.Appointment
	if err := c.BindJSON(&appointmentData); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	region, exists := c.Get("region")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized region"})
		return
	}
	regionStr, ok := region.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid region type"})
		return
	}
	db, err := database.GetDBForRegion(regionStr)
	var doctor database.Doctors
	if err := db.Where("doctor_id = ?", appointmentData.DoctorID).First(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctor information"})
		return
	}

	// Fetch user details
	var user database.Patients
	if err := db.Where("patient_id = ?", appointmentData.PatientID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	appointment := database.Appointment{
		PatientID:       appointmentData.PatientID,
		DoctorID:        appointmentData.DoctorID,
		AppointmentDate: appointmentData.AppointmentDate,
		AppointmentTime: appointmentData.AppointmentTime,
		Description:     appointmentData.Description,
	}

	// Generate the hospital username based on HospitalID, HospitalName, and AdminID
	appointments, err := json.Marshal(appointment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal hospital admin data to JSON"})
		return
	}

	var errKafka error
	switch regionStr {
	case "north":
		// Send to North region's Kafka topic
		errKafka = kafkaManager.SendHospitalRegistrationMessage(regionStr, "appointment_reg", string(appointments))
	case "south":
		// Send to South region's Kafka topic
		errKafka = kafkaManager.SendHospitalRegistrationMessage(regionStr, "appointment_reg", string(appointments))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid region: %s", region)})
		return
	}

	// Check if there was an error sending the message to Kafka
	if errKafka != nil {
		log.Printf("Failed to send hospital registration data to Kafka: %v", errKafka)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Kafka"})
		return
	}
	// Send appointment email to the user

	c.JSON(http.StatusCreated, gin.H{"message": "Appointment created successfully", "appointment_id": appointment.AppointmentID})
}
