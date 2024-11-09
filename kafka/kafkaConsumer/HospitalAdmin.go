// package kafkaconsumer

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/IBM/sarama"
// 	"github.com/adityjoshi/avinyaa/database" // Assuming your DB models are here
// )

// // ConsumerGroupHandler implements the sarama.ConsumerGroupHandler interface.
// type ConsumerGroupHandler struct {
// 	Region string
// }

// // Setup is called when the consumer group is initialized.
// func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
// 	log.Printf("Consumer group setup for region: %s", h.Region)
// 	return nil
// }

// // Cleanup is called when the consumer group is closed.
// func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
// 	log.Printf("Consumer group cleanup for region: %s", h.Region)
// 	return nil
// }

// // ConsumeClaim is called to process messages from the assigned Kafka partition.
// func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
// 	// Process each message
// 	for message := range claim.Messages() {
// 		log.Printf("Consumed message from region %s: %s", h.Region, string(message.Value))

// 		// Example: Write data to the DB based on region
// 		err := processMessage(h.Region, message.Value)
// 		if err != nil {
// 			log.Printf("Error processing message: %v", err)
// 			continue
// 		}

// 		// Mark message as processed
// 		sess.MarkMessage(message, "")
// 	}

// 	return nil
// }

// // processMessage processes the message based on region (North or South).
// func processMessage(region string, message []byte) error {
// 	// Depending on the region, write to the correct DB
// 	if region == "north" {
// 		return saveToNorthDB(message)
// 	} else if region == "south" {
// 		return saveToSouthDB(message)
// 	}
// 	return fmt.Errorf("invalid region: %s", region)
// }

// // saveToNorthDB saves the message to the North database.
// func saveToNorthDB(message []byte) error {
// 	// Here, you would decode the message and save it to NorthDB
// 	// Example: Assume message is a JSON object for a Hospital Admin registration

// 	// Placeholder code: Replace with actual logic to save to NorthDB
// 	var admin database.HospitalAdmin
// 	if err := json.Unmarshal(message, &admin); err != nil {
// 		return fmt.Errorf("failed to unmarshal message for NorthDB: %v", err)
// 	}

// 	if err := database.NorthDB.Create(&admin).Error; err != nil {
// 		return fmt.Errorf("failed to save to NorthDB: %v", err)
// 	}

// 	log.Printf("Saved to NorthDB: %v", admin)
// 	return nil
// }

// // saveToSouthDB saves the message to the South database.
// func saveToSouthDB(message []byte) error {
// 	// Placeholder code: Replace with actual logic to save to SouthDB
// 	var admin database.HospitalAdmin
// 	if err := json.Unmarshal(message, &admin); err != nil {
// 		return fmt.Errorf("failed to unmarshal message for SouthDB: %v", err)
// 	}

// 	if err := database.SouthDB.Create(&admin).Error; err != nil {
// 		return fmt.Errorf("failed to save to SouthDB: %v", err)
// 	}

// 	log.Printf("Saved to SouthDB: %v", admin)
// 	return nil
// }

package kafkaconsumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/adityjoshi/avinyaa/database"
)

// ConsumerGroupHandler implements the sarama.ConsumerGroupHandler interface.
type ConsumerGroupHandler struct {
	Region string
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group setup for region: %s", h.Region)
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group cleanup for region: %s", h.Region)
	return nil
}

// ConsumeClaim is called to process messages from the assigned Kafka partition.
func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Process each message
	for message := range claim.Messages() {
		log.Printf("Consumed message from region %s: %s", h.Region, string(message.Value))

		// Example: Write data to the DB based on region
		err := processMessage(h.Region, message.Value)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		// Mark message as processed
		sess.MarkMessage(message, "")
	}

	return nil
}

// processMessage processes the message based on region (North or South).
func processMessage(region string, message []byte) error {
	if region == "north" {
		return saveToNorthDB(message)
	} else if region == "south" {
		return saveToSouthDB(message)
	}
	return fmt.Errorf("invalid region: %s", region)
}

func saveToNorthDB(message []byte) error {
	var admin database.HospitalAdmin
	if err := json.Unmarshal(message, &admin); err != nil {
		return fmt.Errorf("failed to unmarshal message for NorthDB: %v", err)
	}

	if err := database.NorthDB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to save to NorthDB: %v", err)
	}

	log.Printf("Saved to NorthDB: %v", admin)
	return nil
}

func saveToSouthDB(message []byte) error {
	var admin database.HospitalAdmin
	if err := json.Unmarshal(message, &admin); err != nil {
		return fmt.Errorf("failed to unmarshal message for SouthDB: %v", err)
	}

	if err := database.SouthDB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to save to SouthDB: %v", err)
	}

	log.Printf("Saved to SouthDB: %v", admin)
	return nil
}
