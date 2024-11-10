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

package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/adityjoshi/avinyaa/database"
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

	case "hospital_registrations":
		// Process hospital updates messages
		log.Printf("Processing hospital_registration message: %s", string(msg.Value))
		// Add your logic for processing hospital_updates messages here
	default:
		// Handle any other topics or log an error if the topic is not recognized
		log.Printf("Received message from unknown topic: %s", topic)
		// Add your default logic here
	}
	return nil
}
