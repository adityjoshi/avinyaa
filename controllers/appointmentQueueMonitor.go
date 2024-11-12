//package controllers

// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"fmt"
// // 	"log"
// // 	"time"

// // 	"github.com/adityjoshi/avinyaa/database" // Import your database package to access the DB connections
// // 	"github.com/adityjoshi/avinyaa/utils"

// // 	"gorm.io/gorm"
// // )

// // func CheckAppointmentsQueue() {
// // 	// Monitor the Redis queues for all hospitals and departments
// // 	departments := []string{"Cardiology", "Neurology", "Orthopedics", "Pediatrics", "General"} // Add your departments here

// // 	// Check appointment queues for each department in each region (North, South, etc.)
// // 	regions := []string{"north", "south", "east", "west"}
// // 	for _, region := range regions {
// // 		for _, department := range departments {
// // 			// Start the process for each region and department
// // 			go monitorQueueForDepartment(region, department)
// // 		}
// // 	}

// // 	// Keep the function running (or use a wait group to control the flow of execution)
// // 	select {}
// // }

// // // monitorQueueForDepartment will be responsible for checking the appointment queue for a specific department in a region.
// // func monitorQueueForDepartment(region, department string) {
// // 	var db *gorm.DB
// // 	var redisKey string

// // 	// Select the appropriate database and Redis key based on the region
// // 	switch region {
// // 	case "north":
// // 		db = database.NorthDB
// // 	case "south":
// // 		db = database.SouthDB
// // 	case "east":
// // 		db = database.DB // Assuming `DB` is the default database for east region
// // 	case "west":
// // 		db = database.DB // Assuming `DB` is the default database for west region
// // 	}

// // 	redisKey = fmt.Sprintf("appointments:%s:%s", region, department)

// // 	// Infinite loop to keep checking the Redis queue for new appointments
// // 	for {
// // 		// Fetch appointments from Redis queue
// // 		result, err := database.RedisClient.LRange(context.Background(), redisKey, 0, -1).Result()
// // 		if err != nil {
// // 			log.Printf("Error fetching appointments from Redis queue: %v", err)
// // 			continue
// // 		}

// // 		for _, item := range result {
// // 			// Process each appointment (decode the JSON message)
// // 			var appointment database.Appointment
// // 			if err := json.Unmarshal([]byte(item), &appointment); err != nil {
// // 				log.Printf("Error unmarshalling appointment data: %v", err)
// // 				continue
// // 			}

// // 			// Check if this appointment is eligible for notification
// // 			checkAndNotifyAppointment(db, appointment, redisKey)
// // 		}

// // 		// Sleep for a short time to avoid continuous querying of Redis
// // 		time.Sleep(10 * time.Second)
// // 	}
// // }

// // // checkAndNotifyAppointment checks the appointment queue, calculates the number of people ahead, and sends a notification
// // func checkAndNotifyAppointment(db *gorm.DB, appointment database.Appointment, redisKey string) {
// // 	// Count the number of appointments in the queue for the same department
// // 	count, err := database.RedisClient.LLen(context.Background(), redisKey).Result()
// // 	if err != nil {
// // 		log.Printf("Error counting appointments in Redis queue: %v", err)
// // 		return
// // 	}

// // 	// Get the position of the current patient in the queue
// // 	position := count - 1 // 0-indexed, so the last item has the position `count-1`

// // 	// Notify the patient when their appointment is close (e.g., 10 people ahead)
// // 	if position <= 10 {
// // 		notifyPatientOfUpcomingAppointment(appointment)
// // 	}

// // 	// Additional logic to send notifications when the patient is nearing their appointment time can go here
// // }

// // // notifyPatientOfUpcomingAppointment sends a notification to the patient when their appointment is near
// // func notifyPatientOfUpcomingAppointment(appointment database.Appointment) {
// // 	// Fetch the patient details based on the PatientID
// // 	var patient database.Patients
// // 	if err := database.DB.Where("patient_id = ?", appointment.PatientID).First(&patient).Error; err != nil {
// // 		log.Printf("Error fetching patient details for PatientID %s: %v", appointment.PatientID, err)
// // 		return
// // 	}
// // 	var doctor database.Doctors
// // 	if err := database.DB.Where("doctor_id = ?", appointment.DoctorID).First(&doctor).Error; err != nil {
// // 		log.Printf("Error fetching doctor details for DoctorID %s: %v", appointment.DoctorID, err)
// // 		return
// // 	}
// // 	// Send a notification (e.g., email, SMS, etc.)
// // 	// For example, print a message to the console or send an email
// // 	err := utils.SendAppointmentEmail(patient.Email, doctor.FullName, appointment.AppointmentDate.Format("2006-01-02"), appointment.AppointmentTime.Format("15:04"), time.December.String())
// // 	if err != nil {
// // 		log.Fatalf("Error sending appointment email: %v", err)
// // 	}

// // 	fmt.Printf("Patient %s, your appointment with Dr. %s is approaching!\n", patient.FullName, doctor.FullName)
// // }

// package controllers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/adityjoshi/avinyaa/database" // Import your database package
// 	"github.com/adityjoshi/avinyaa/utils"    // Utility functions like email notifications
// 	"gorm.io/gorm"
// )

// // CheckAppointmentsQueue monitors the Redis queues for all hospitals and departments
// func CheckAppointmentsQueue() {
// 	// Monitor the Redis queues for all hospitals and departments
// 	departments := []string{"Cardiology", "Neurology", "Orthopedics", "Pediatrics", "General"} // Add your departments here

// 	// Check appointment queues for each department in each region (North, South, etc.)
// 	regions := []string{"north", "south", "east", "west"}
// 	for _, region := range regions {
// 		for _, department := range departments {
// 			// Start the process for each region and department
// 			go monitorQueueForDepartment(region, department)
// 		}
// 	}

// 	// Keep the function running (or use a wait group to control the flow of execution)
// 	select {}
// }

// // monitorQueueForDepartment checks the appointment queue for a specific department in a region
// func monitorQueueForDepartment(region, department string) {
// 	// Define a dynamic Redis key based on region, department, and hospital ID
// 	redisKey := fmt.Sprintf("appointments:%s:%s", region, department)

// 	// Infinite loop to keep checking the Redis queue for new appointments
// 	for {
// 		// Fetch appointments from Redis queue
// 		result, err := database.RedisClient.LRange(context.Background(), redisKey, 0, -1).Result()
// 		if err != nil {
// 			log.Printf("Error fetching appointments from Redis queue: %v", err)
// 			continue
// 		}

// 		// Process each appointment
// 		for _, item := range result {
// 			var appointment database.Appointment
// 			if err := json.Unmarshal([]byte(item), &appointment); err != nil {
// 				log.Printf("Error unmarshalling appointment data: %v", err)
// 				continue
// 			}

// 			// Check and notify about the appointment
// 			checkAndNotifyAppointment(region, appointment, redisKey)
// 		}

// 		// Sleep for a short time to avoid continuous querying of Redis
// 		time.Sleep(10 * time.Second)
// 	}
// }

// // checkAndNotifyAppointment processes the appointment and checks if it's close to being next in line
// func checkAndNotifyAppointment(region string, appointment database.Appointment, redisKey string) {
// 	// Select the correct database based on the region
// 	var db *gorm.DB
// 	switch region {
// 	case "north":
// 		db = database.NorthDB
// 	case "south":
// 		db = database.SouthDB
// 	case "east":
// 		db = database.DB
// 	case "west":
// 		db = database.DB
// 	default:
// 		log.Printf("Unknown region: %v", region)
// 		return
// 	}

// 	// Count the number of appointments in the queue for the same department
// 	count, err := database.RedisClient.LLen(context.Background(), redisKey).Result()
// 	if err != nil {
// 		log.Printf("Error counting appointments in Redis queue: %v", err)
// 		return
// 	}

// 	// Get the position of the current patient in the queue
// 	position := count - 1 // 0-indexed, so the last item has the position `count-1`

// 	// Notify the patient if their appointment is close (e.g., 10 people ahead)
// 	if position <= 1 {
// 		notifyPatientOfUpcomingAppointment(db, appointment)
// 	}
// }

// // notifyPatientOfUpcomingAppointment sends an email notification to the patient
// func notifyPatientOfUpcomingAppointment(db *gorm.DB, appointment database.Appointment) {
// 	// Fetch the patient details based on the PatientID
// 	var patient database.Patients
// 	if err := db.Where("patient_id = ?", appointment.PatientID).First(&patient).Error; err != nil {
// 		log.Printf("Error fetching patient details for PatientID %s: %v", appointment.PatientID, err)
// 		return
// 	}

// 	// Fetch the doctor details based on the DoctorID
// 	var doctor database.Doctors
// 	if err := db.Where("doctor_id = ?", appointment.DoctorID).First(&doctor).Error; err != nil {
// 		log.Printf("Error fetching doctor details for DoctorID %s: %v", appointment.DoctorID, err)
// 		return
// 	}

// 	// Send an email notification to the patient
// 	err := utils.SendAppointmentEmail(patient.Email, doctor.FullName, appointment.AppointmentDate.Format("2006-01-02"), appointment.AppointmentTime.Format("15:04"), "December")
// 	if err != nil {
// 		log.Printf("Error sending appointment email: %v", err)
// 	}

// 	// Log the notification
// 	fmt.Printf("Patient %s, your appointment with Dr. %s is approaching!\n", patient.FullName, doctor.FullName)
// }

// working model below

// package controllers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/adityjoshi/avinyaa/database" // Import your database package
// 	"github.com/adityjoshi/avinyaa/utils"    // Utility functions like email notifications
// 	"gorm.io/gorm"
// )

// // CheckAppointmentsQueue monitors the Redis queues for all hospitals and departments
// func CheckAppointmentsQueue() {
// 	// Monitor the Redis queues for all hospitals and departments
// 	departments := []string{"Cardiology", "Neurology", "Orthopedics", "Pediatrics", "General"} // Add your departments here

// 	// Check appointment queues for each department in each region (North, South, etc.)
// 	regions := []string{"North", "South", "East", "West"}
// 	hospitalIDs := []string{"1", "2", "3"} // Example hospital IDs, modify as needed

// 	// Check for each hospital and department in each region
// 	for _, region := range regions {
// 		for _, department := range departments {
// 			for _, hospitalID := range hospitalIDs {
// 				// Start the process for each region, department, and hospital
// 				go monitorQueueForDepartment(region, hospitalID, department)
// 			}
// 		}
// 	}

// 	// Keep the function running (or use a wait group to control the flow of execution)
// 	select {}
// }

// // monitorQueueForDepartment checks the appointment queue for a specific department, hospital, and region
// func monitorQueueForDepartment(region, hospitalID, department string) {
// 	redisKey := fmt.Sprintf("appointments:%s:%s:%s", region, hospitalID, department)

// 	for {
// 		result, err := database.RedisClient.LRange(context.Background(), redisKey, 0, -1).Result()
// 		if err != nil {
// 			log.Printf("Error fetching appointments from Redis queue: %v", err)
// 			continue
// 		}

// 		for _, item := range result {
// 			var appointment database.Appointment
// 			if err := json.Unmarshal([]byte(item), &appointment); err != nil {
// 				log.Printf("Error unmarshalling appointment data: %v", err)
// 				continue
// 			}

// 			// Check if the appointment is marked as "done" in the database
// 			if appointment.IsDone {
// 				continue // Skip the appointment if it is already done
// 			}

// 			// Process the appointment if it's not marked as done
// 			checkAndNotifyAppointment(region, hospitalID, department, appointment, redisKey)
// 		}

// 		time.Sleep(10 * time.Second)
// 	}
// }

// // checkAndNotifyAppointment processes the appointment and checks if it's close to being next in line
// func checkAndNotifyAppointment(region, hospitalID, department string, appointment database.Appointment, redisKey string) {
// 	// Select the correct database based on the region
// 	var db *gorm.DB
// 	switch region {
// 	case "North":
// 		db = database.NorthDB
// 	case "South":
// 		db = database.SouthDB
// 	case "East":
// 		db = database.DB
// 	case "West":
// 		db = database.DB
// 	default:
// 		log.Printf("Unknown region: %v", region)
// 		return
// 	}

// 	// Count the number of appointments in the queue for the same department, hospital, and region
// 	count, err := database.RedisClient.LLen(context.Background(), redisKey).Result()
// 	if err != nil {
// 		log.Printf("Error counting appointments in Redis queue: %v", err)
// 		return
// 	}

// 	// Log the count of appointments to verify the queue length
// 	log.Printf("Number of appointments in the queue for %s:%s:%s: %d", region, hospitalID, department, count)

// 	// Get the position of the current patient in the queue
// 	position := count - 1 // 0-indexed, so the last item has the position `count-1`
// 	log.Printf("Position of the current patient in the queue: %d", position)

// 	// Notify the patient if their appointment is close (e.g., 10 people ahead)
// 	if position <= 10 {
// 		notifyPatientOfUpcomingAppointment(db, &appointment)
// 	}
// }

// // notifyPatientOfUpcomingAppointment sends an email notification to the patient
// // func notifyPatientOfUpcomingAppointment(db *gorm.DB, appointment database.Appointment) {
// // 	// Fetch the patient details based on the PatientID
// // 	var patient database.Patients
// // 	if err := db.Where("patient_id = ?", appointment.PatientID).First(&patient).Error; err != nil {
// // 		log.Printf("Error fetching patient details for PatientID %s: %v", appointment.PatientID, err)
// // 		return
// // 	}

// // 	// Fetch the doctor details based on the DoctorID
// // 	var doctor database.Doctors
// // 	if err := db.Where("doctor_id = ?", appointment.DoctorID).First(&doctor).Error; err != nil {
// // 		log.Printf("Error fetching doctor details for DoctorID %s: %v", appointment.DoctorID, err)
// // 		return
// // 	}

// // 	// Log before sending the email to check if it is being triggered
// // 	fmt.Printf("Sending email to patient %s: %s, with doctor %s: %s\n", patient.FullName, patient.Email, doctor.FullName, doctor.FullName)

// // 	// Send an email notification to the patient
// // 	err := utils.SendAppointmentEmail(patient.Email, doctor.FullName, appointment.AppointmentDate.Format("2006-01-02"), appointment.AppointmentTime.Format("15:04"), "December")
// // 	if err != nil {
// // 		log.Printf("Error sending appointment email: %v", err)
// // 	}

// //		// Log the notification
// //		fmt.Printf("Patient %s, your appointment with Dr. %s is approaching!\n", patient.FullName, doctor.FullName)
// //	}
// func notifyPatientOfUpcomingAppointment(db *gorm.DB, appointment *database.Appointment) {
// 	// Skip notification if the appointment is marked as done
// 	if appointment.IsDone {
// 		return
// 	}

// 	var patient database.Patients
// 	if err := db.Where("patient_id = ?", appointment.PatientID).First(&patient).Error; err != nil {
// 		log.Printf("Error fetching patient details for PatientID %s: %v", appointment.PatientID, err)
// 		return
// 	}

// 	var doctor database.Doctors
// 	if err := db.Where("doctor_id = ?", appointment.DoctorID).First(&doctor).Error; err != nil {
// 		log.Printf("Error fetching doctor details for DoctorID %s: %v", appointment.DoctorID, err)
// 		return
// 	}

// 	// Send email to the patient
// 	fmt.Printf("Sending email to patient %s: %s, with doctor %s: %s\n", patient.FullName, patient.Email, doctor.FullName, doctor.FullName)

// 	err := utils.SendAppointmentEmail(patient.Email, doctor.FullName, appointment.AppointmentDate.Format("2006-01-02"), appointment.AppointmentTime.Format("15:04"), "December")
// 	if err != nil {
// 		log.Printf("Error sending appointment email: %v", err)
// 	}

// 	// Log that the email has been sent
// 	fmt.Printf("Patient %s, your appointment with Dr. %s is approaching!\n", patient.FullName, doctor.FullName)

// 	// Mark the appointment as "done" after sending the email
// 	appointment.IsDone = true

// 	// Save the updated appointment status in the database
// 	if err := db.Save(appointment).Error; err != nil {
// 		log.Printf("Error updating appointment status to done for PatientID %d: %v", appointment.PatientID, err)
// 	}
// }

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adityjoshi/avinyaa/database"
	"github.com/adityjoshi/avinyaa/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// CheckAppointmentsQueue monitors the Redis queues for all hospitals and departments
func CheckAppointmentsQueue() {
	// Monitor the Redis queues for all hospitals and departments
	departments := []string{"Cardiology", "Neurology", "Orthopedics", "Pediatrics", "General"} // Add your departments here

	// Check appointment queues for each department in each region (North, South, etc.)
	regions := []string{"North", "South", "East", "West"}
	hospitalIDs := []string{"1", "2", "3"}

	// Check for each hospital and department in each region
	for _, region := range regions {
		for _, department := range departments {
			for _, hospitalID := range hospitalIDs {
				// Start the process for each region, department, and hospital
				go monitorQueueForDepartment(region, hospitalID, department)
			}
		}
	}

	// Keep the function running (or use a wait group to control the flow of execution)
	select {}
}

// monitorQueueForDepartment checks the appointment queue for a specific department, hospital, and region
func monitorQueueForDepartment(region, hospitalID, department string) {
	redisKey := fmt.Sprintf("appointments:%s:%s:%s", region, hospitalID, department)
	var appointment database.Appointment
	for {
		result, err := database.RedisClient.LRange(context.Background(), redisKey, 0, -1).Result()
		if err != nil {
			log.Printf("Error fetching appointments from Redis queue: %v", err)
			continue
		}

		for _, item := range result {

			if err := json.Unmarshal([]byte(item), &appointment); err != nil {
				log.Printf("Error unmarshalling appointment data: %v", err)
				continue
			}

			// Check if the appointment is marked as "done" in the database
			if appointment.Appointed {
				// Remove the appointment from the Redis queue if it is done
				err := database.RedisClient.LRem(context.Background(), redisKey, 0, item).Err()
				if err != nil {
					log.Printf("Error removing appointment from Redis queue: %v", err)
				} else {
					log.Printf("Removed appointment %v from Redis queue", appointment.AppointmentID)
				}
				continue // Skip the appointment if it is already done
			}

			// Process the appointment if it's not marked as done
			checkAndNotifyAppointment(region, hospitalID, department, appointment, redisKey)
		}

		time.Sleep(10 * time.Second)
	}
}

// checkAndNotifyAppointment processes the appointment and checks if it's close to being next in line
func checkAndNotifyAppointment(region, hospitalID, department string, appointment database.Appointment, redisKey string) {
	// Select the correct database based on the region
	var db *gorm.DB
	switch region {
	case "North":
		db = database.NorthDB
	case "South":
		db = database.SouthDB
	case "East":
		db = database.DB
	case "West":
		db = database.DB
	default:
		log.Printf("Unknown region: %v", region)
		return
	}

	// Count the number of appointments in the queue for the same department, hospital, and region
	count, err := database.RedisClient.LLen(context.Background(), redisKey).Result()
	if err != nil {
		log.Printf("Error counting appointments in Redis queue: %v", err)
		return
	}

	// Log the count of appointments to verify the queue length
	log.Printf("Number of appointments in the queue for %s:%s:%s: %d", region, hospitalID, department, count)

	// Get the position of the current patient in the queue
	position := count - 1 // 0-indexed, so the last item has the position `count-1`
	log.Printf("Position of the current patient in the queue: %d", position)

	// Notify the patient if their appointment is close (e.g., 10 people ahead)
	if position <= 10 {
		notifyPatientOfUpcomingAppointment(db, &appointment, region)
	}

	// After notifying or processing, remove the appointment from the queue if it's done
	if appointment.Appointed {
		// Construct the correct appointment data to match Redis format
		appointmentData := fmt.Sprintf("%d:%d", appointment.AppointmentID, appointment.PatientID)

		// Ensure the Redis key is correctly constructed
		redisKey := fmt.Sprintf("appointments:%s:%s:%s", region, hospitalID, department)

		// Remove the appointment from the Redis queue
		removed, err := database.RedisClient.LRem(context.Background(), redisKey, 0, appointmentData).Result()
		if err != nil {
			log.Printf("Error removing appointment from Redis queue: %v", err)
		} else if removed > 0 {
			log.Printf("Successfully removed appointment %d from Redis queue", appointment.AppointmentID)
		} else {
			log.Printf("No appointment found in the queue to remove for AppointmentID %d", appointment.AppointmentID)
		}
	}
}

// isEmailSent checks if an email has already been sent for the given patient and region
func isEmailSent(redisKey, region, patientID string) bool {
	emailKey := fmt.Sprintf("%s:%s:emailSent:%s", redisKey, region, patientID)
	status, err := database.RedisClient.Get(context.Background(), emailKey).Result()
	if err == redis.Nil {
		return false // Email hasn't been sent yet
	} else if err != nil {
		log.Printf("Error checking Redis for email sent status: %v", err)
		return false
	}

	return status == "true"
}

// markEmailAsSent marks the email as sent for a given patient and region
func markEmailAsSent(redisKey, region, patientID string) {
	emailKey := fmt.Sprintf("%s:%s:emailSent:%s", redisKey, region, patientID)
	if err := database.RedisClient.Set(context.Background(), emailKey, "true", 24*time.Hour).Err(); err != nil {
		log.Printf("Error marking email as sent in Redis: %v", err)
	}
}

// notifyPatientOfUpcomingAppointment sends an email notification to the patient
func notifyPatientOfUpcomingAppointment(db *gorm.DB, appointment *database.Appointment, region string) {
	// Skip notification if the appointment is marked as done
	if appointment.IsDone {
		return
	}

	// Check if email has already been sent for this appointment
	redisKey := fmt.Sprintf("appointments:%s", region) // Redis key format
	if isEmailSent(redisKey, region, strconv.Itoa(int(appointment.PatientID))) {
		log.Printf("Email already sent for patient %s in region %s. Skipping notification.", appointment.PatientID, region)
		return
	}

	// Fetch patient details
	var patient database.Patients
	if err := db.Where("patient_id = ?", appointment.PatientID).First(&patient).Error; err != nil {
		log.Printf("Error fetching patient details for PatientID %s: %v", appointment.PatientID, err)
		return
	}

	// Fetch doctor details
	var doctor database.Doctors
	if err := db.Where("doctor_id = ?", appointment.DoctorID).First(&doctor).Error; err != nil {
		log.Printf("Error fetching doctor details for DoctorID %s: %v", appointment.DoctorID, err)
		return
	}

	// Send email to the patient
	fmt.Printf("Sending email to patient %s: %s, with doctor %s: %s\n", patient.FullName, patient.Email, doctor.FullName, doctor.FullName)

	err := utils.SendAppointmentEmail(patient.Email, doctor.FullName, appointment.AppointmentDate.Format("2006-01-02"), appointment.AppointmentTime.Format("15:04"), "December")
	if err != nil {
		log.Printf("Error sending appointment email: %v", err)
	}

	// Log that the email has been sent
	fmt.Printf("Patient %s, your appointment with Dr. %s is approaching!\n", patient.FullName, doctor.FullName)

	// Mark the appointment as "done" after sending the email
	appointment.IsDone = true

	// Save the updated appointment status in the database
	if err := db.Save(appointment).Error; err != nil {
		log.Printf("Error updating appointment status to done for PatientID %d: %v", appointment.PatientID, err)
	}

	// After sending the email, mark it as sent in Redis
	markEmailAsSent(redisKey, region, strconv.Itoa(int(appointment.PatientID)))
}
