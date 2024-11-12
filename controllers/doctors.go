package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/adityjoshi/avinyaa/database"
	"github.com/adityjoshi/avinyaa/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// func RegisterDoctor(c *gin.Context) {
// 	var doctorData struct {
// 		FullName      string              `json:"full_name"`
// 		Description   string              `json:"description"`
// 		ContactNumber string              `json:"contact_number"`
// 		Email         string              `json:"email"`
// 		AdminID       uint                `json:"admin_id"`
// 		Department    database.Department `json:"department"`
// 	}

// 	if err := c.BindJSON(&doctorData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Ensure AdminID is included in the JSON payload
// 	if doctorData.AdminID == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin ID is required"})
// 		return
// 	}

// 	// Find the hospital associated with the admin
// 	var hospital database.Hospitals
// 	if err := database.DB.Where("admin_id = ?", doctorData.AdminID).First(&hospital).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found for the admin"})
// 		return
// 	}

// 	// Set HospitalID and HospitalName in doctor data
// 	doctor := database.Doctors{
// 		FullName:      doctorData.FullName,
// 		Description:   doctorData.Description,
// 		ContactNumber: doctorData.ContactNumber,
// 		Email:         doctorData.Email,
// 		HospitalID:    hospital.HospitalId,   // Correctly set HospitalID from fetched hospital
// 		Hospital:      hospital.HospitalName, // Set HospitalName
// 		Department:    doctorData.Department,
// 	}

// 	// Generate username
// 	doctor.Username = generateDoctorUsername(doctor.HospitalID, hospital.HospitalName, doctor.FullName)

// 	// Save doctor data
// 	if err := database.DB.Create(&doctor).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register doctor"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Doctor registered successfully", "hospital_name": hospital.HospitalName})
// }

func RegisterDoctor(c *gin.Context) {

	var doctorData database.Doctors

	if err := c.BindJSON(&doctorData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Extract AdminID from JWT claims
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
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

	// Find the hospital associated with the admin
	var hospital database.Hospitals
	db, err := database.GetDBForRegion(regionStr)
	if err = db.Where("admin_id = ?", adminIDUint).First(&hospital).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found for the admin"})
		return
	}

	password := generatePassword(doctorData.FullName, regionStr)
	doctorData.Password = password
	fmt.Print(doctorData.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctorData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)

	}
	doctorData.Password = string(hashedPassword)
	// Set HospitalID and HospitalName in doctor data
	doctor := database.Doctors{
		FullName:      doctorData.FullName,
		Description:   doctorData.Description,
		ContactNumber: doctorData.ContactNumber,
		Email:         doctorData.Email,
		HospitalID:    hospital.HospitalId,   // Correctly set HospitalID from fetched hospital
		Hospital:      hospital.HospitalName, // Set HospitalName
		Department:    doctorData.Department,
		Password:      doctorData.Password,
		Region:        regionStr,
	}

	// Generate username
	doctor.Username = generateDoctorUsername(doctor.HospitalID, doctor.FullName)

	// Save doctor data
	if err := db.Create(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register doctor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor registered successfully", "hospital_name": hospital.HospitalName})
}
func generatePassword(fullName, region string) string {
	// For simplicity, we combine full name and username to generate a password
	fmt.Sprintf("%s%s", fullName, region)
	return fmt.Sprintf("%s%s", fullName, region)
}

// Helper function to generate doctor username
func generateDoctorUsername(hospitalID uint, doctorFullName string) string {
	// Remove spaces from hospital name and doctor full name

	doctorFullName = strings.ReplaceAll(doctorFullName, " ", "")
	// Construct username
	return fmt.Sprintf("%d%s%s", hospitalID, doctorFullName)
}

func GetDoctor(c *gin.Context) {
	doctorID := c.Param("doctor_id")
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

	var doctor database.Doctors
	db, err := database.GetDBForRegion(regionStr)
	if err = db.First(&doctor, doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	var hospital database.Hospitals
	if err := db.First(&hospital, doctor.HospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found for the doctor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"doctor_id":      doctor.DoctorID,
		"full_name":      doctor.FullName,
		"description":    doctor.Description,
		"contact_number": doctor.ContactNumber,
		"email":          doctor.Email,
		"hospital_id":    doctor.HospitalID,
		"hospital_name":  hospital.HospitalName,
		"department":     doctor.Department,
		"username":       doctor.Username,
		"region":         doctor.Region,
	})
}

func DoctorLogin(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Region   string `json:"region"`
	}
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var doctor database.Doctors
	db, err := database.GetDBForRegion(loginRequest.Region)
	if err := db.Where("email = ?", loginRequest.Email).First(&doctor).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(doctor.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}
	otp, err := GenerateAndSendOTP(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate or send OTP" + otp})
		return
	}
	// func GenerateJwt(userID uint, userType, role, region string)
	token, err := utils.GenerateDoctorJwt(doctor.DoctorID, "Doctor", "Doctor", string(doctor.Department), loginRequest.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with message to enter OTP
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email. Please verify the OTP.", "token": token, "region": loginRequest.Region})
}

// MarkAppointmentDone marks an appointment as done and removes it from the Redis queue
// func MarkAppointmentAsDone(c *gin.Context) {
// 	// Retrieve the JWT token from the Authorization header
// 	tokenString := c.GetHeader("Authorization")
// 	if tokenString == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
// 		return
// 	}

// 	// Decode the JWT token
// 	claims, err := utils.DecodeJwt(tokenString)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
// 		return
// 	}

// 	// Extract user_type from JWT claims
// 	userType, ok := claims["user_type"].(string)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user type"})
// 		return
// 	}

// 	// Ensure that only doctors are allowed to mark appointments as done
// 	if userType != "Doctor" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only doctors can mark appointments as done."})
// 		return
// 	}

// 	// Extract and validate doctor-specific information from JWT claims
// 	doctorID, doctorIDExists := claims["doctor_id"].(float64) // Change to float64
// 	department, departmentExists := claims["department"].(string)
// 	region, regionExists := claims["region"].(string)

// 	if !doctorIDExists || !departmentExists || !regionExists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor credentials are incomplete in token"})
// 		return
// 	}

// 	doctorIDUint := uint(doctorID) // Convert from float64 to uint

// 	// Set doctor info in the context for further use if needed
// 	c.Set("doctor_id", doctorIDUint)
// 	c.Set("department", department)
// 	c.Set("region", region)

// 	// Get the appointment ID from the request parameters
// 	appointmentIDParam := c.Param("appointment_id")
// 	appointmentID, err := strconv.ParseUint(appointmentIDParam, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
// 		return
// 	}

// 	// Select the correct database based on the region from JWT claims
// 	db, err := database.GetDBForRegion(region)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to regional database"})
// 		return
// 	}

// 	var hospital database.Doctors
// 	if err := db.Where("doctor_id = ? AND region = ?", doctorIDUint, region).First(&hospital).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hospital ID for doctor"})
// 		return
// 	}

// 	hospitalID := hospital.HospitalID

// 	// Retrieve the appointment from the database
// 	var appointment database.Appointment
// 	if err := db.Where("appointment_id = ? AND doctor_id = ?", appointmentID, doctorIDUint).First(&appointment).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		}
// 		return
// 	}

// 	// Mark the appointment as done
// 	appointment.Appointed = true

// 	if err := db.Save(&appointment).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment status"})
// 		return
// 	}

// 	// Optionally, remove the patient from the Redis queue and notify other services
// 	// hospitalIDStr := strconv.FormatUint(uint64(hospitalID), 10)
// 	// redisKey := fmt.Sprintf("appointments:%s:%d:%s", region, hospitalIDStr, department)
// 	// if _, err := database.RedisClient.LRem(c, redisKey, 0, appointmentID).Result(); err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove appointment from queue"})
// 	// 	return
// 	// }

// 	redisKey := fmt.Sprintf("appointments:%s:%d:%s", region, hospitalID, department)
// 	appointmentData := fmt.Sprintf("%d:%d", appointmentID, appointment.PatientID)

// 	// Remove the appointment from the Redis queue
// 	if _, err := database.RedisClient.LRem(c, redisKey, 0, appointmentData).Result(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove appointment from queue"})
// 		return
// 	}

// 	// Respond to indicate the appointment was successfully marked as done
// 	c.JSON(http.StatusOK, gin.H{"message": "Appointment marked as done successfully"})
// }

func MarkAppointmentAsDone(c *gin.Context) {
	// Retrieve the JWT token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}

	// Decode the JWT token
	claims, err := utils.DecodeJwt(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
		return
	}

	// Extract user_type from JWT claims
	userType, ok := claims["user_type"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user type"})
		return
	}

	// Ensure that only doctors are allowed to mark appointments as done
	if userType != "Doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Only doctors can mark appointments as done."})
		return
	}

	// Extract and validate doctor-specific information from JWT claims
	doctorID, doctorIDExists := claims["doctor_id"].(float64)
	department, departmentExists := claims["department"].(string)
	region, regionExists := claims["region"].(string)

	if !doctorIDExists || !departmentExists || !regionExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Doctor credentials are incomplete in token"})
		return
	}

	doctorIDUint := uint(doctorID) // Convert from float64 to uint

	// Set doctor info in the context for further use if needed
	c.Set("doctor_id", doctorIDUint)
	c.Set("department", department)
	c.Set("region", region)

	// Get the appointment ID from the request parameters
	appointmentIDParam := c.Param("appointment_id")
	appointmentID, err := strconv.ParseUint(appointmentIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	// Select the correct database based on the region from JWT claims
	db, err := database.GetDBForRegion(region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to regional database"})
		return
	}

	var hospital database.Doctors
	if err := db.Where("doctor_id = ? AND region = ?", doctorIDUint, region).First(&hospital).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hospital ID for doctor"})
		return
	}

	hospitalID := hospital.HospitalID

	var appointment database.Appointment
	if err := db.Where("appointment_id = ? AND doctor_id = ?", appointmentID, doctorIDUint).First(&appointment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Mark the appointment as done
	appointment.Appointed = true

	if err := db.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment status"})
		return
	}
	hospitalIDStr := strconv.Itoa(int(hospitalID))

	fmt.Printf("Hospital ID as string: %s\n", hospitalIDStr)

	redisKey := fmt.Sprintf("appointments:%s:%s:%s", region, hospitalIDStr, department)
	appointmentData := fmt.Sprintf("%d:%d", appointmentID, appointment.PatientID)

	if _, err := database.RedisClient.LRem(c, redisKey, 0, appointmentData).Result(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove appointment from queue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment marked as done successfully"})
}
