package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/VictorMont03/AutoTracking/golang-simulator/internal"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("Starting simulator...")

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	freightTopic := os.Getenv("KAFKA_FREIGHT_TOPIC")
	simulationTopic := os.Getenv("KAFKA_SIMULATION_TOPIC")
	routeTopic := os.Getenv("KAFKA_ROUTE_TOPIC")
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")

	mongoConnection, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	freightService := internal.NewFreightService()
	routeService := internal.NewRouteService(mongoConnection, freightService)

	chDriverMoved := make(chan *internal.DriverMovedEvent)

	kafkaBroker := os.Getenv("KAFKA_BROKER")

	freightWrite := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    freightTopic,
		Balancer: &kafka.LeastBytes{},
	}

	simulatorWriter := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    simulationTopic,
		Balancer: &kafka.LeastBytes{},
	}

	listener := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   routeTopic,
		GroupID: kafkaGroupId,
	})

	hub := internal.NewEventHub(routeService, mongoConnection, chDriverMoved, freightWrite, simulatorWriter)

	fmt.Println("Listening for messages...")
	for {
		msg, err := listener.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Error reading message: ", err)
			continue
		}

		go func(msg []byte) {
			err = hub.HandleEvent(msg)
			if err != nil {
				fmt.Println("Error handling event: ", err)

			}
		}(msg.Value)
	}
}
