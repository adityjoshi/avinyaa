package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/adityjoshi/avinyaa/database"
	kafkamanager "github.com/adityjoshi/avinyaa/kafka/kafkaManager"
	"github.com/adityjoshi/avinyaa/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterHospitalAdmin(c *gin.Context) {
	// Retrieve KafkaManager from the context
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

	var admin database.HospitalAdmin

	// Bind incoming JSON request to the HospitalAdmin struct
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default the Usertype if not provided
	if admin.Usertype == "" {
		admin.Usertype = "Admin"
	}

	// Check if the user already exists in the database
	// var existingUser database.HospitalAdmin
	var existingUser database.HospitalAdmin
	db, err := database.GetDBForRegion(admin.Region)
	if err := db.Where("email = ?", admin.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	admin.Password = string(hashedPassword)

	// Get the region from the request body (ensure region is part of the request)
	region := admin.Region // Assuming the region is passed as form data
	if region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Region is required"})
		return
	}

	// Prepare the message to send to Kafka (just using the admin data here)
	//message := fmt.Sprintf("Admin ID: %s, Name: %s, Email: %s, Usertype: %s", admin.AdminID, admin.FullName, admin.Email, admin.Password, admin.ContactNumber, admin.Region, admin.Usertype)
	adminMessage, err := json.Marshal(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal hospital admin data to JSON"})
		return
	}
	// Send the registration message to Kafka based on the region
	var errKafka error
	switch region {
	case "north":
		// Send to North region's Kafka topic (you provide the topic name)
		errKafka = kafkaManager.SendUserRegistrationMessage(region, "hospital_admin", string(adminMessage))
	case "south":
		// Send to South region's Kafka topic (you provide the topic name)
		errKafka = kafkaManager.SendUserRegistrationMessage(region, "hospital_admin", string(adminMessage))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid region: %s", region)})
		return
	}

	// Check if there was an error sending the message to Kafka
	if errKafka != nil {
		log.Printf("Failed to send registration data to Kafka: %v", errKafka)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Kafka"})
		return
	}

	// Respond with success if both Kafka message sending and database insertion were successful
	c.JSON(http.StatusCreated, gin.H{
		"message":              "Hospital admin registered successfully",
		"admin_id":             admin.AdminID,
		"kafka_message_status": "Message sent to Kafka successfully",
	})
}

func AdminLogin(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Region   string `json:"region"`
	}
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin database.HospitalAdmin
	db, err := database.GetDBForRegion(loginRequest.Region)
	if err := db.Where("email = ?", loginRequest.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}
	otp, err := GenerateAndSendOTP(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate or send OTP" + otp})
		return
	}
	// func GenerateJwt(userID uint, userType, role, region string)
	token, err := utils.GenerateJwt(admin.AdminID, "Admin", string(admin.Usertype), admin.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with message to enter OTP
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email. Please verify the OTP.", "token": token, "region": admin.Region})
}

func VerifyAdminOTP(c *gin.Context) {
	var otpRequest struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.BindJSON(&otpRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get region from context and assert its type
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

	// Verify the OTP
	isValid, err := VerifyOtp(otpRequest.Email, otpRequest.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying OTP"})
		return
	}
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	// Retrieve the appropriate database based on region
	db, err := database.GetDBForRegion(regionStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to regional database"})
		return
	}

	// Retrieve user information after OTP verification
	var user database.HospitalAdmin
	if err := db.Where("email = ?", otpRequest.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Set OTP verification status in Redis
	redisClient := database.GetRedisClient()
	err = redisClient.Set(context.Background(), "otp_verified:"+strconv.Itoa(int(user.AdminID)), "verified", 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error setting OTP verification status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"loggedin": "success", "region": regionStr})
}

func RegisterHospital(c *gin.Context) {
	var hospital database.Hospitals
	if err := c.BindJSON(&hospital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}
	hospital.AdminID = adminIDUint
	region, exists := c.Get("region")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Region not specified"})
		return
	}
	regionStr, ok := region.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid region type"})
		return
	}
	hospital.Region = regionStr

	// Retrieve the appropriate database based on region
	db, err := database.GetDBForRegion(regionStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to regional database"})
		return
	}

	var latestHospital database.Hospitals
	if err := db.Order("hospital_id DESC").First(&latestHospital).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve latest hospital"})
			return
		}
	}

	// Generate the hospital username based on HospitalID, HospitalName, and AdminID
	hospitalRegistration, err := json.Marshal(hospital)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal hospital admin data to JSON"})
		return
	}

	var errKafka error
	switch regionStr {
	case "north":
		// Send to North region's Kafka topic
		errKafka = kafkaManager.SendHospitalRegistrationMessage(regionStr, "hospital_registration", string(hospitalRegistration))
	case "south":
		// Send to South region's Kafka topic
		errKafka = kafkaManager.SendHospitalRegistrationMessage(regionStr, "hospital_registration", string(hospitalRegistration))
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

	c.JSON(http.StatusCreated, gin.H{"message": "Hospital created successfully", "region": regionStr})
}

func GetHospital(c *gin.Context) {
	hospitalID := c.Param("hospital_id")

	var hospital database.Hospitals
	if err := database.DB.First(&hospital, hospitalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		return
	}

	c.JSON(http.StatusOK, hospital)
}

func RegisterStaff(c *gin.Context) {
	var staff database.HospitalStaff
	if err := c.BindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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

	// Get admin ID from JWT
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}

	// Verify admin's hospital authorization

	region, exists := c.Get("region")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Region not specified"})
		return
	}
	regionStr, ok := region.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid region type"})
		return
	}

	hospitalID, err := verifyAdminHospital(adminIDUint, regionStr)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin not authorized to register staff"})
		return
	}

	db, err := database.GetDBForRegion(regionStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error connecting to regional database"})
		return
	}
	// Retrieve hospital details
	var hospital database.Hospitals
	if err := db.Where("hospital_id = ?", hospitalID).First(&hospital).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hospital details"})
		return
	}
	staff.Region = regionStr
	staff.HospitalID = hospitalID
	staff.Region = regionStr
	staff.HospitalName = hospital.HospitalName

	staffRegister, err := json.Marshal(staff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal hospital staff data to JSON"})
		return
	}
	var errKafka error
	switch regionStr {
	case "north":
		// Send to North region's Kafka topic
		errKafka = kafkaManager.SendHospitalStaffRegisterMessage(regionStr, "hospital_staff", string(staffRegister))
	case "south":
		// Send to South region's Kafka topic
		errKafka = kafkaManager.SendHospitalStaffRegisterMessage(regionStr, "hospital_staff", string(staffRegister))
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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Staff created successfully",
		"region":  regionStr,
		// Optionally return the generated password
	})
}

func AddBedType(c *gin.Context) {
	var bedsCount database.BedsCount

	// Parse the JSON request body into the bedsCount struct
	if err := c.BindJSON(&bedsCount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get admin ID from JWT
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}

	// Verify the admin's hospital
	hospitalID, err := verifyAdminHospital(adminIDUint, "")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin not authorized to add beds for this hospital"})
		return
	}

	// Set the hospital ID to the bedsCount object
	bedsCount.HospitalID = hospitalID

	// Check if the bed type already exists for the hospital
	var existingBedType database.BedsCount
	if err := database.DB.Where("hospital_id = ? AND type_name = ?", hospitalID, bedsCount.TypeName).First(&existingBedType).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Bed type already exists for this hospital"})
		return
	}

	// Save the new bed type and total beds
	if err := database.DB.Create(&bedsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add bed type"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Bed type added successfully",
		"bed_type_id": bedsCount.ID,
		"type_name":   bedsCount.TypeName,
		"total_beds":  bedsCount.TotalBeds,
		"hospital_id": bedsCount.HospitalID,
	})
}

func UpdateTotalBeds(c *gin.Context) {
	var bedData struct {
		TypeName  string `json:"type_name"`
		TotalBeds int    `json:"total_beds"` // Number of beds to add or remove
		Action    string `json:"action"`     // Action: "add" or "remove"
	}

	// Parse the JSON request body into the bedData struct
	if err := c.BindJSON(&bedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get admin ID from JWT
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}

	// Verify the admin's hospital
	hospitalID, err := verifyAdminHospital(adminIDUint, "")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin not authorized to update beds for this hospital"})
		return
	}

	// Find the bed type for the given hospital
	var bedType database.BedsCount
	if err := database.DB.Where("hospital_id = ? AND type_name = ?", hospitalID, bedData.TypeName).First(&bedType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bed type not found for this hospital"})
		return
	}

	// Perform the add or remove action
	switch bedData.Action {
	case "add":
		// Add beds
		previousTotalBeds := bedType.TotalBeds
		bedType.TotalBeds += uint(bedData.TotalBeds)

		// Add rooms for the newly added beds
		for i := previousTotalBeds + 1; i <= bedType.TotalBeds; i++ {
			roomNumber := fmt.Sprintf("%s%d", strings.ToLower(bedData.TypeName), i)
			newRoom := database.Room{
				HospitalID: hospitalID,
				BedType:    bedData.TypeName,
				RoomNumber: roomNumber,
				IsOccupied: false,
			}
			if err := database.DB.Create(&newRoom).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rooms"})
				return
			}
		}

	case "remove":
		// Remove beds (ensure total beds don't go below zero)
		if int(bedType.TotalBeds)-bedData.TotalBeds < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove more beds than available"})
			return
		}

		// Remove unoccupied rooms first
		var unoccupiedRooms []database.Room
		if err := database.DB.Where("hospital_id = ? AND bed_type = ? AND is_occupied = ?", hospitalID, bedData.TypeName, false).Order("room_number desc").Limit(bedData.TotalBeds).Find(&unoccupiedRooms).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find unoccupied rooms"})
			return
		}

		// Ensure there are enough unoccupied rooms to remove
		if len(unoccupiedRooms) < bedData.TotalBeds {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough unoccupied rooms to remove"})
			return
		}

		// Delete the unoccupied rooms
		for _, room := range unoccupiedRooms {
			if err := database.DB.Delete(&room).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
				return
			}
		}

		// Update total beds
		bedType.TotalBeds -= uint(bedData.TotalBeds)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Use 'add' or 'remove'"})
		return
	}

	// Save the updated bed count
	if err := database.DB.Save(&bedType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bed count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Total beds updated successfully",
		"type_name":  bedType.TypeName,
		"total_beds": bedType.TotalBeds,
	})
}

func GetTotalBeds(c *gin.Context) {
	// Get admin ID from JWT
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid admin ID"})
		return
	}

	// Verify the admin's hospital
	hospitalID, err := verifyAdminHospital(adminIDUint, "")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin not authorized to view beds for this hospital"})
		return
	}

	// Query the BedsCount for the given hospital
	var beds []database.BedsCount
	if err := database.DB.Where("hospital_id = ?", hospitalID).Find(&beds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bed information"})
		return
	}

	// If no beds are found for the hospital
	if len(beds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No bed data found for this hospital"})
		return
	}

	// Create a response that lists all the bed types with their total, available, and occupied counts
	var bedDetails []gin.H
	for _, bed := range beds {
		bedDetails = append(bedDetails, gin.H{
			"type_name":  bed.TypeName,
			"total_beds": bed.TotalBeds,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"hospital_id": hospitalID,
		"bed_details": bedDetails,
	})
}

func verifyAdminHospital(adminID uint, region string) (uint, error) {
	var admin database.HospitalAdmin
	db, err := database.GetDBForRegion(region)
	if err != nil {
		return 0, err // Return error if database connection failed
	}
	if err := db.Where("admin_id = ?", adminID).First(&admin).Error; err != nil {
		return 0, err
	}

	var hospital database.Hospitals
	if err := db.Where("admin_id = ?", adminID).First(&hospital).Error; err != nil {
		return 0, err
	}

	return hospital.HospitalId, nil
}
