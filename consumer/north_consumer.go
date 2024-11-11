package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/adityjoshi/avinyaa/database"
	"golang.org/x/crypto/bcrypt"
)

type NorthConsumer struct {
	Consumer sarama.Consumer
	Topics   []string
}

// NewNorthConsumer initializes a new Kafka consumer for the north region (multiple topics)
func NewNorthConsumer(broker string, topics []string) (*NorthConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create the consumer
	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Printf("Error creating consumer: %v", err)
		return nil, fmt.Errorf("error creating consumer: %v", err)
	}

	log.Printf("Kafka consumer created successfully, subscribing to topics: %v", topics)

	// Return a NorthConsumer instance with the list of topics
	return &NorthConsumer{Consumer: consumer, Topics: topics}, nil
}

// Listen starts the consumer and listens for messages on the specified topics
func (nc *NorthConsumer) Listen() {
	defer func() {
		if err := nc.Consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v\n", err)
		}
	}()

	// Create a consumer for each topic, listening to partition 0 for each one
	for _, topic := range nc.Topics {
		partitionConsumer, err := nc.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to start consumer for partition 0 of topic %s: %v", topic, err)
		}
		defer partitionConsumer.Close()

		// Log that the consumer has started listening
		log.Printf("Consumer is now listening to topic: %s from partition 0\n", topic)

		// Start a goroutine for each topic to consume messages concurrently
		go nc.consumeMessages(partitionConsumer)
	}

	// Block indefinitely (since we're listening to multiple topics concurrently)
	select {}
}

// consumeMessages handles message consumption for each topic
func (nc *NorthConsumer) consumeMessages(partitionConsumer sarama.PartitionConsumer) {
	for msg := range partitionConsumer.Messages() {
		// Log received message and metadata
		log.Printf("Received message from topic %s: %s\n", msg.Topic, string(msg.Value))
		log.Printf("Message metadata - Partition: %d, Offset: %d\n", msg.Partition, msg.Offset)

		// Process the message (e.g., save data to the database or trigger further actions)
		if err := processMessage(msg.Topic, msg); err != nil {
			log.Printf("Error processing message: %v", err)
		} else {
			log.Printf("Message processed successfully")
		}
	}
}

// processMessage is a placeholder function for processing the received Kafka message
func processMessage(topic string, msg *sarama.ConsumerMessage) error {
	if database.NorthDB == nil {
		log.Fatal("NorthDB is not initialized!")
		return fmt.Errorf("NorthDB is not initialized")
	}

	log.Printf("Processing message: %s \n", string(msg.Value))
	switch topic {
	case "hospital_admin":
		// Process hospital admin messages
		log.Printf("Processing hospital_admin message: %s", string(msg.Value))

		var admin database.HospitalAdmin
		if err := json.Unmarshal(msg.Value, &admin); err != nil {
			log.Printf("Failed to unmarshal hospital_admin message: %v", err)
			return err
		}
		if err := database.NorthDB.Create(&admin); err != nil {
			log.Printf("Failed to save hospital_admin data: %v", err.Error)
			return fmt.Errorf("Failed to write to the DB", err.Error)
		}
		// Add your logic for processing hospital_admin messages here

	case "hospital_registration":
		log.Printf("Processing hospital_registration message: %s", string(msg.Value))

		var hospital database.Hospitals
		if err := json.Unmarshal(msg.Value, &hospital); err != nil {
			log.Printf("Error unmarshalling hospital data: %v", err)
			return err

		}
		hospital.Username = fmt.Sprintf("DEL%d", hospital.HospitalId)

		if err := database.NorthDB.Create(&hospital).Error; err != nil {
			log.Printf("Error creating hospital in database: %v", err)
			return fmt.Errorf(err.Error())
		}

	case "hospital_staff":
		log.Printf("Processing hospital_staff registration message: %s", string(msg.Value))

		var staff database.HospitalStaff
		err := json.Unmarshal(msg.Value, &staff)
		if err != nil {
			log.Printf("Failed to unmarshal hospital staff data: %v", err)
			return fmt.Errorf("failed to unmarshal staff data: %v", err)
		}

		password := generatePassword(staff.FullName, staff.Region)
		staff.Password = password

		// Hash the staff password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			return fmt.Errorf("failed to hash password: %v", err)
		}
		staff.Password = string(hashedPassword)
		staff.Username = fmt.Sprintf("%s%s", staff.ContactNumber, strings.ReplaceAll(strings.ToLower(staff.FullName), " ", ""))

		// Save the staff to the database
		if err := database.NorthDB.Create(&staff).Error; err != nil {
			log.Printf("Failed to save staff data to database: %v", err)
			return fmt.Errorf("failed to save staff data to database: %v", err)
		}

		log.Printf("Staff member created successfully with ID: %d", staff.StaffID)
		return nil
	case "patient_registration":
		log.Printf("Processing patient_registration message: %s", string(msg.Value))

		var patients database.Patients
		if err := json.Unmarshal(msg.Value, &patients); err != nil {
			log.Printf("Error unmarshalling patients data: %v", err)
			return err

		}

		if err := database.NorthDB.Create(&patients).Error; err != nil {
			log.Printf("Error creating patients in database: %v", err)
			return fmt.Errorf(err.Error())
		}
	case "patient_Admit":
		log.Printf("Processing patient_registration message: %s", string(msg.Value))

		var patient_Admit database.PatientBeds
		if err := json.Unmarshal(msg.Value, &patient_Admit); err != nil {
			log.Printf("Error unmarshalling patient_admit data: %v", err)
			return err

		}

		if err := database.NorthDB.Create(&patient_Admit).Error; err != nil {
			log.Printf("Error creating patients in database: %v", err)
			return fmt.Errorf(err.Error())
		}
		var availableRoom database.Room
		availableRoom.IsOccupied = true
		if err := database.NorthDB.Save(&availableRoom).Error; err != nil {
			log.Printf("Error saving patient_admit data: %v", err)
			return err
		}

	default:
		// Handle any other topics or log an error if the topic is not recognized
		log.Printf("Received message from unknown topic: %s", topic)
		// Add your default logic here
	}
	return nil
}

func generatePassword(fullName, username string) string {
	// For simplicity, we combine full name and username to generate a password
	return fmt.Sprintf("%s%s", fullName, username)
}

// // consumer/north_consumer.go

// package consumer

// import (
// 	"fmt"
// 	"log"

// 	"github.com/IBM/sarama"
// )

// type NorthConsumer struct {
// 	Consumer sarama.Consumer
// 	Topic    string
// }

// // NewNorthConsumer initializes a new Kafka consumer for the north region
// func NewNorthConsumer(broker, topic string) (*NorthConsumer, error) {
// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true

// 	consumer, err := sarama.NewConsumer([]string{broker}, config)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating consumer: %v", err)
// 	}

// 	return &NorthConsumer{Consumer: consumer, Topic: topic}, nil
// }

// // Listen starts the consumer and listens for messages on the specified topic
// func (nc *NorthConsumer) Listen() {
// 	defer func() {
// 		if err := nc.Consumer.Close(); err != nil {
// 			log.Printf("Error closing consumer: %v\n", err)
// 		}
// 	}()

// 	partitionConsumer, err := nc.Consumer.ConsumePartition(nc.Topic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		log.Fatalf("Failed to start consumer for partition: %v", err)
// 	}
// 	defer partitionConsumer.Close()

// 	fmt.Printf("Listening to messages from north region on topic: %s\n", nc.Topic)

// 	for msg := range partitionConsumer.Messages() {
// 		fmt.Printf("Received message: %s\n", string(msg.Value))
// 		// Here you can add logic to save data to the north region database
// 	}
// }

// consumer/north_consumer.go
// consumer/north_consumer.go
// consumer/north_consumer.go

// package consumer

// import (
// 	"fmt"
// 	"log"

// 	"github.com/IBM/sarama"
// )

// type NorthConsumer struct {
// 	Consumer sarama.Consumer
// 	Topic    string
// }

// // NewNorthConsumer initializes a new Kafka consumer for the north region
// func NewNorthConsumer(broker, topic string) (*NorthConsumer, error) {
// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true

// 	consumer, err := sarama.NewConsumer([]string{broker}, config)
// 	if err != nil {
// 		log.Printf("Error creating consumer: %v", err)
// 		return nil, fmt.Errorf("error creating consumer: %v", err)
// 	}

// 	log.Printf("Kafka consumer created successfully, subscribing to topic: %s", topic)

// 	return &NorthConsumer{Consumer: consumer, Topic: topic}, nil
// }

// // Listen starts the consumer and listens for messages on the specified topic
// func (nc *NorthConsumer) Listen() {
// 	defer func() {
// 		if err := nc.Consumer.Close(); err != nil {
// 			log.Printf("Error closing consumer: %v\n", err)
// 		}
// 	}()

// 	// Consume messages from partition 0, starting from the newest message
// 	partitionConsumer, err := nc.Consumer.ConsumePartition(nc.Topic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		log.Fatalf("Failed to start consumer for partition: %v", err)
// 	}
// 	defer partitionConsumer.Close()

// 	// Log that the consumer has started listening
// 	log.Printf("Consumer is now listening to topic: %s from partition 0\n", nc.Topic)

// 	// Continuous loop to receive messages
// 	for msg := range partitionConsumer.Messages() {
// 		// Log received message and metadata
// 		log.Printf("Received message from topic %s: %s\n", msg.Topic, string(msg.Value))
// 		log.Printf("Message metadata - Partition: %d, Offset: %d\n", msg.Partition, msg.Offset)

// 		// Process the message (e.g., save data to the database or trigger further actions)
// 		if err := processMessage(msg); err != nil {
// 			log.Printf("Error processing message: %v", err)
// 		} else {
// 			log.Printf("Message processed successfully")
// 		}
// 	}
// }

// // processMessage is a placeholder function for processing the received Kafka message
// func processMessage(msg *sarama.ConsumerMessage) error {
// 	// Implement the logic for handling the message, like saving to the database
// 	log.Printf("Processing message: %s", string(msg.Value))
// 	return nil
// }
