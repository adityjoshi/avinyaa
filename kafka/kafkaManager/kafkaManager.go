// // package kafka

// // import (
// // 	"fmt"

// // 	"github.com/adityjoshi/avinyaa/kafka"
// // )

// // // KafkaManager is responsible for managing Kafka producers and sending messages to topics
// // type KafkaManager struct {
// // 	northProducer *kafka.NorthProducer
// // 	southProducer *kafka.SouthProducer
// // }

// // // NewKafkaManager initializes and returns a KafkaManager instance
// // func NewKafkaManager(northBrokers, southBrokers []string) (*KafkaManager, error) {
// // 	// Initialize the North producer
// // 	northProducer, err := kafka.NewNorthProducer(northBrokers)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("Error initializing North producer: %w", err)
// // 	}

// // 	// Initialize the South producer
// // 	southProducer, err := kafka.NewSouthProducer(southBrokers)
// // 	if err != nil {
// // 		return nil, fmt.Errorf("Error initializing South producer: %w", err)
// // 	}

// // 	// Return the KafkaManager instance with both producers
// // 	return &KafkaManager{
// // 		northProducer: northProducer,
// // 		southProducer: southProducer,
// // 	}, nil
// // }

// // // SendUserRegistrationMessage sends the user registration data to the appropriate Kafka topic based on the region
// // func (km *KafkaManager) SendUserRegistrationMessage(region, topic, message string) error {
// // 	var err error

// // 	// Determine which producer to use based on the region
// // 	switch region {
// // 	case "north":
// // 		// Use the North producer
// // 		err = km.northProducer.SendMessage(topic, message)
// // 	case "south":
// // 		// Use the South producer
// // 		err = km.southProducer.SendMessage(topic, message)
// // 	default:
// // 		return fmt.Errorf("invalid region: %s", region)
// // 	}

// // 	// Return any errors from sending the message
// // 	if err != nil {
// // 		return fmt.Errorf("failed to send message to Kafka: %w", err)
// // 	}
// // 	return nil
// // }

// package kafkamanager

// import (
// 	"fmt"
// 	"log"

// 	"github.com/adityjoshi/avinyaa/kafka" // Correct path to producer
// )

// // KafkaManager is responsible for managing Kafka producers and sending messages to topics
// type KafkaManager struct {
// 	northProducer *kafka.NorthProducer
// 	southProducer *kafka.SouthProducer
// }

// // NewKafkaManager initializes and returns a KafkaManager instance
// func NewKafkaManager(northBrokers, southBrokers []string) (*KafkaManager, error) {
// 	// Initialize the North producer
// 	northProducer, err := kafka.NewNorthProducer(northBrokers)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error initializing North producer: %w", err)
// 	}

// 	// Initialize the South producer
// 	southProducer, err := kafka.NewSouthProducer(southBrokers)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error initializing South producer: %w", err)
// 	}

// 	// Return the KafkaManager instance with both producers
// 	return &KafkaManager{
// 		northProducer: northProducer,
// 		southProducer: southProducer,
// 	}, nil
// }

// // SendUserRegistrationMessage sends the user registration data to the appropriate Kafka topic based on the region and registration type
// func (km *KafkaManager) SendUserRegistrationMessage(region, topic, message string) error {
// 	var err error

// 	// Check if the region is valid and determine the correct producer
// 	switch region {
// 	case "north":
// 		// Send to the North region's producer
// 		err = km.northProducer.SendMessage(topic, message)
// 	case "south":
// 		// Send to the South region's producer
// 		err = km.southProducer.SendMessage(topic, message)
// 	default:
// 		return fmt.Errorf("invalid region: %s", region)
// 	}

// 	// Return any errors from sending the message
// 	if err != nil {
// 		return fmt.Errorf("failed to send message to Kafka: %w", err)
// 	}

// 	log.Printf("Message sent to topic %s in %s region", topic, region)
// 	return nil
// }

package kafkamanager

import (
	"fmt"
	"log"

	"github.com/adityjoshi/avinyaa/kafka" // Correct path to producer
)

// KafkaManager is responsible for managing Kafka producers and sending messages to topics
type KafkaManager struct {
	northProducer *kafka.NorthProducer
	southProducer *kafka.SouthProducer
}

// NewKafkaManager initializes and returns a KafkaManager instance
func NewKafkaManager(northBrokers, southBrokers []string) (*KafkaManager, error) {
	// Initialize the North producer
	northProducer, err := kafka.NewNorthProducer(northBrokers)
	if err != nil {
		return nil, fmt.Errorf("Error initializing North producer: %w", err)
	}

	// Initialize the South producer
	southProducer, err := kafka.NewSouthProducer(southBrokers)
	if err != nil {
		return nil, fmt.Errorf("Error initializing South producer: %w", err)
	}

	// Return the KafkaManager instance with both producers
	return &KafkaManager{
		northProducer: northProducer,
		southProducer: southProducer,
	}, nil
}

// SendUserRegistrationMessage sends the user registration data to the appropriate Kafka topic based on the region and registration type
func (km *KafkaManager) SendUserRegistrationMessage(region, topic, message string) error {
	var err error

	// Check if the region is valid and determine the correct producer
	switch region {
	case "north":
		// Log for debugging
		log.Printf("Sending message to North region, topic: %s", topic)
		err = km.northProducer.SendMessage(topic, message)
	case "south":
		// Log for debugging
		log.Printf("Sending message to South region, topic: %s", topic)
		err = km.southProducer.SendMessage(topic, message)
	default:
		// Return an error if the region is invalid
		return fmt.Errorf("invalid region: %s", region)
	}

	// Return any errors from sending the message
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka topic %s in %s region: %w", topic, region, err)
	}

	// Log successful message sending
	log.Printf("Message successfully sent to topic %s in %s region", topic, region)
	return nil
}

func (km *KafkaManager) SendHospitalRegistrationMessage(region, topic, message string) error {
	var err error

	// Check if the region is valid and determine the correct producer
	switch region {
	case "north":
		// Log for debugging
		log.Printf("Sending message to North region, topic: %s", topic)
		err = km.northProducer.SendMessage(topic, message)
	case "south":
		// Log for debugging
		log.Printf("Sending message to South region, topic: %s", topic)
		err = km.southProducer.SendMessage(topic, message)
	default:
		// Return an error if the region is invalid
		return fmt.Errorf("invalid region: %s", region)
	}

	// Return any errors from sending the message
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka topic %s in %s region: %w", topic, region, err)
	}

	// Log successful message sending
	log.Printf("Message successfully sent to topic %s in %s region", topic, region)
	return nil
}

func (km *KafkaManager) SendHospitalStaffRegisterMessage(region, topic, message string) error {
	var err error

	// Check if the region is valid and determine the correct producer
	switch region {
	case "north":
		// Log for debugging
		log.Printf("Sending message to North region, topic: %s", topic)
		err = km.northProducer.SendMessage(topic, message)
	case "south":
		// Log for debugging
		log.Printf("Sending message to South region, topic: %s", topic)
		err = km.southProducer.SendMessage(topic, message)
	default:
		// Return an error if the region is invalid
		return fmt.Errorf("invalid region: %s", region)
	}

	// Return any errors from sending the message
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka topic %s in %s region: %w", topic, region, err)
	}

	// Log successful message sending
	log.Printf("Message successfully sent to topic %s in %s region", topic, region)
	return nil
}
