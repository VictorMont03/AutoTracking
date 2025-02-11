package internal

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Route struct {
	ID           string      `json:"id" bson:"_id"`
	Distance     int         `json:"distance" bson:"distance"`
	Directions   []Direction `json:"directions" bson:"directions"`
	FreightPrice float64     `json:"freightPrice" bson:"freightPrice"`
}

type Direction struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

type RouteService struct {
	mongo          *mongo.Client
	freightService *FreightService
}

func NewRouteService(mongo *mongo.Client, freightService *FreightService) *RouteService {
	return &RouteService{
		mongo:          mongo,
		freightService: freightService,
	}
}

func (rs *RouteService) CreateRoute(route Route) (Route, error) {
	fmt.Println(`Route Service | Started`)

	route.FreightPrice = rs.freightService.CalculateFreightPrice(route.Distance)

	update := bson.M{
		"$set": bson.M{
			"distance":     route.Distance,
			"directions":   route.Directions,
			"freightPrice": route.FreightPrice,
		},
	}

	filter := bson.M{"_id": route.ID}
	opts := options.Update().SetUpsert(true)

	fmt.Println(`Route Service | Updating route`)

	_, err := rs.mongo.Database("routes").Collection("routes").UpdateOne(nil, filter, update, opts)

	if err != nil {
		return Route{}, err
	}

	fmt.Println(`Route Service | Route updated`)

	return route, nil
}

func (rs *RouteService) GetRoute(id string) (Route, error) {
	var route Route

	filter := bson.M{"_id": id}

	err := rs.mongo.Database("routes").Collection("routes").FindOne(nil, filter).Decode(&route)

	if err != nil {
		return Route{}, err
	}

	return route, nil
}
