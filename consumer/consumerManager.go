// consumer/consumerFile.go

package consumer

import (
	"fmt"
	"log"
)

func StartConsumer(region string) {
	kafkaBroker := "localhost:9092"
	var topic = []string{
		"hospital_admin",
		"hospital_registration",
		"hospital_staff",
		"patient_registration",
		"patient_Admit",
		"patient_admission",

		// Add other topics as necessary
	}
	switch region {
	case "north":
		//topic := "hospital_admin"
		northConsumer, err := NewNorthConsumer(kafkaBroker, topic)
		if err != nil {
			log.Fatalf("Failed to create north consumer: %v", err)
		}
		northConsumer.Listen()

	case "south":
		// Similar setup for the south region consumer
		fmt.Println("Starting south consumer...")

	case "east":
		// Similar setup for the east region consumer
		fmt.Println("Starting east consumer...")

	case "west":
		// Similar setup for the west region consumer
		fmt.Println("Starting west consumer...")

	default:
		fmt.Println("Unknown region:", region)
	}
}
