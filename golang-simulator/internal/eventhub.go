package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event interface {
	EventName() string
}

type EventHub struct {
	RouteService    *RouteService
	mongoClient     *mongo.Client
	chDriverMoved   chan *DriverMovedEvent
	freightWriter   *kafka.Writer
	simulatorWriter *kafka.Writer
}

func NewEventHub(routeService *RouteService, mongoClient *mongo.Client, chDriverMoved chan *DriverMovedEvent, freightWriter *kafka.Writer, simulatorWriter *kafka.Writer) *EventHub {
	return &EventHub{
		RouteService:    routeService,
		mongoClient:     mongoClient,
		chDriverMoved:   chDriverMoved,
		freightWriter:   freightWriter,
		simulatorWriter: simulatorWriter,
	}
}

func (eh *EventHub) HandleEvent(msg []byte) error {
	//
	fmt.Println("Handling event | ", string(msg))

	var baseEvent struct {
		EventName string `json:"eventName"`
	}

	err := json.Unmarshal(msg, &baseEvent)

	if err != nil {
		return err
	}

	switch baseEvent.EventName {
	case "RouteCreated":
		var event RouteCreatedEvent
		err := json.Unmarshal(msg, &event)

		if err != nil {
			fmt.Println("Error unmarshalling RouteCreatedEvent")
			return err
		}

		freightCalculatedEvent, err := RouteCreatedHandler(&event, eh.RouteService)

		if err != nil {
			fmt.Println("Error handling RouteCreatedEvent")
			return err
		}

		decodeFreight, err := json.Marshal(freightCalculatedEvent)

		if err != nil {
			fmt.Println("Error marshalling FreightCalculatedEvent")
			return err
		}

		err = eh.PublishMessage(baseEvent.EventName, decodeFreight)

		if err != nil {
			fmt.Println("Error publishing FreightCalculatedEvent")
			return err
		}

	case "DeliveryStarted":
		var event DeliveryStartedEvent
		err := json.Unmarshal(msg, &event)

		if err != nil {
			fmt.Println("Error unmarshalling DeliveryStartedEvent")
			return err
		}

		err = DeliveryStartedHandler(&event, eh.RouteService, eh.chDriverMoved)

		fmt.Println("Delivery Started Handler called")

		if err != nil {
			fmt.Println("Error handling DeliveryStarted")
			return err
		}

		fmt.Println("Initiating go routine to process DriverMovedEvent")

		return eh.ProcessRoutine(eh.chDriverMoved)
	default:
		return errors.New("event not found")
	}

	return nil
}

func (eh *EventHub) ProcessRoutine(events chan *DriverMovedEvent) error {
	go func() {
		for {
			select {
			case movedEvent := <-events:
				message, err := json.Marshal(movedEvent)

				if err != nil {
					fmt.Println("Error marshalling event")
				}

				err = eh.PublishMessage(movedEvent.EventName, message)

				if err != nil {
					fmt.Println("Error publishing event")
				}
			case <-time.After(500 * time.Millisecond):
				fmt.Println("Timeout waiting for event")
				return
			}
		}
	}()

	return nil
}

func (eh *EventHub) PublishMessage(eventName string, message []byte) error {
	//
	var err error

	switch eventName {
	case "RouteCreated":
		err = eh.freightWriter.WriteMessages(context.Background(), kafka.Message{Value: message})

	case "DriverMoved":
		err = eh.simulatorWriter.WriteMessages(context.Background(), kafka.Message{Value: message})

	}

	if err != nil {
		fmt.Println("Error writing message to Kafka")
		fmt.Println(err)
		return err
	}

	return nil

}
