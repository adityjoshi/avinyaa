// package kafkaconsumer

// import (
// 	"log"

// 	"github.com/IBM/sarama"
// )

// // StartConsumer initializes and starts Kafka consumers for different regions and topics.
// func StartConsumer(brokers []string, groupID, region string) {
// 	// Set up the Kafka consumer group config
// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true

// 	// Create the consumer group
// 	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
// 	if err != nil {
// 		log.Fatalf("Error creating consumer group: %v", err)
// 	}

// 	// Create the consumer handler for the region (North or South)
// 	handler := &ConsumerGroupHandler{Region: region}

// 	// Topics for which we will start consuming messages
// 	// List of topics (e.g. "hospitaladmin", "appointments", etc.)
// 	topics := []string{"hospitaladmin", "appointments", "payment", "staff", "admissions"}

// 	// Start consuming messages for each topic
// 	go func() {
// 		for _, topic := range topics {
// 			// Each topic may have its own consumer group handler or shared handler
// 			err := consumerGroup.Consume(nil, []string{topic}, handler)
// 			if err != nil {
// 				log.Fatalf("Error consuming messages from topic %s: %v", topic, err)
// 			}
// 		}
// 	}()
// }

// package kafkaconsumer

// import (
// 	"log"

// 	"github.com/IBM/sarama"
// )

// // StartConsumer initializes and starts Kafka consumers for different regions and topics, each with its own consumer group.
// func StartConsumer(brokers []string, region string) {
// 	// Set up the Kafka consumer group config
// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true

// 	// Topics and their corresponding consumer group IDs
// 	topicsWithGroups := map[string]string{
// 		"hospitaladmin": "hospitaladmin-consumer-group",
// 		// "appointments":  "appointments-consumer-group",
// 		// "payment":       "payment-consumer-group",
// 		// "staff":         "staff-consumer-group",
// 		// "admissions":    "admissions-consumer-group",
// 	}

// 	// Start consuming messages for each topic with its own consumer group
// 	for topic, groupID := range topicsWithGroups {
// 		// Create a new consumer group for each topic
// 		consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
// 		if err != nil {
// 			log.Fatalf("Error creating consumer group for topic %s: %v", topic, err)
// 		}

// 		// Create the consumer handler for the region (North or South)
// 		handler := &ConsumerGroupHandler{Region: region}

// 		// Start consuming messages for this topic
// 		go func(topic, groupID string) {
// 			for {
// 				err := consumerGroup.Consume(nil, []string{topic}, handler)
// 				if err != nil {
// 					log.Fatalf("Error consuming messages from topic %s: %v", topic, err)
// 				}
// 			}
// 		}(topic, groupID)
// 	}
// }

package kafkaconsumer

import (
	"log"

	"github.com/IBM/sarama"
)

// StartConsumer initializes and starts Kafka consumers for different regions and topics, each with its own consumer group.
func StartConsumer(brokers []string, region string) {
	// Set up the Kafka consumer group config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Topics and their corresponding consumer group IDs
	topicsWithGroups := map[string]string{
		"hospitaladmin": "hospitaladmin-consumer-group",
		// "appointments":  "appointments-consumer-group",
		// "payment":       "payment-consumer-group",
		// "staff":         "staff-consumer-group",
		// "admissions":    "admissions-consumer-group",
	}

	// Start consuming messages for each topic with its own consumer group
	for topic, groupID := range topicsWithGroups {
		// Create a new consumer group for each topic
		consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
		if err != nil {
			log.Fatalf("Error creating consumer group for topic %s: %v", topic, err)
		}

		// Create the consumer handler for the region (North or South)
		handler := &ConsumerGroupHandler{Region: region}

		// Start consuming messages for this topic
		go func(topic, groupID string) {
			for {
				err := consumerGroup.Consume(nil, []string{topic}, handler)
				if err != nil {
					log.Fatalf("Error consuming messages from topic %s: %v", topic, err)
				}
			}
		}(topic, groupID)
	}
}
