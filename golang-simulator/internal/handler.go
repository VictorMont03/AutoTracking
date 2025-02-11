package internal

import (
	"fmt"
	"time"
)

type RouteCreatedEvent struct {
	EventName  string      `json:"eventName"`
	RouteID    string      `json:"id"`
	Distance   int         `json:"distance"`
	Directions []Direction `json:"directions"`
}

type FreightCalculatedEvent struct {
	EventName string  `json:"eventName"`
	RouteID   string  `json:"route_id"`
	Amount    float64 `json:"amount"`
}

type DeliveryStartedEvent struct {
	EventName string `json:"eventName"`
	RouteID   string `json:"route_id"`
}

type DriverMovedEvent struct {
	EventName string  `json:"eventName"`
	RouteID   string  `json:"route_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

func NewRoute(id string, distance int, directions []Direction) *Route {
	return &Route{
		ID:         id,
		Distance:   distance,
		Directions: directions,
	}
}

func NewFreightCalculatedEvent(routeID string, amount float64) *FreightCalculatedEvent {
	return &FreightCalculatedEvent{
		EventName: "FreightCalculated",
		RouteID:   routeID,
		Amount:    amount,
	}
}

func NewDeliveryStartedEvent(routeID string) *DeliveryStartedEvent {
	return &DeliveryStartedEvent{
		EventName: "DeliveryStarted",
		RouteID:   routeID,
	}
}

func NewDriverMovedEvent(routeID string, lat, lng float64) *DriverMovedEvent {
	return &DriverMovedEvent{
		EventName: "DriverMoved",
		RouteID:   routeID,
		Lat:       lat,
		Lng:       lng,
	}
}

func RouteCreatedHandler(event *RouteCreatedEvent, routeService *RouteService) (*FreightCalculatedEvent, error) {
	//
	fmt.Println("Creating route...")

	route := NewRoute(event.RouteID, event.Distance, event.Directions)

	routeCreated, err := routeService.CreateRoute(*route)

	if err != nil {
		return nil, err
	}

	freightCalculatedEvent := NewFreightCalculatedEvent(routeCreated.ID, routeCreated.FreightPrice)

	return freightCalculatedEvent, nil
}

func DeliveryStartedHandler(event *DeliveryStartedEvent, routeService *RouteService, ch chan *DriverMovedEvent) error {
	//
	fmt.Println("Delivery Started Handler initiated...")
	route, err := routeService.GetRoute(event.RouteID)

	if err != nil {
		return err
	}

	go func() {
		for _, direction := range route.Directions {
			dme := NewDriverMovedEvent(route.ID, direction.Lat, direction.Lng)
			ch <- dme
			time.Sleep(time.Second)
		}
	}()

	return nil
}
