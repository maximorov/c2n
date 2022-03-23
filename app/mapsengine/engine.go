package mapsengine

import (
	"context"
	"go.uber.org/zap"
	"helpers/app/location"
	"log"

	"googlemaps.github.io/maps"
)

type (
	Store interface {
		GetLocationFromSearchText(search string) (*location.Location, error)
	}

	MapsStore struct {
		client *maps.Client
	}
)

func NewStore() MapsStore {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBG3pxHBZMc0pb2cRjzMdd7ss0XugNNdf4"))
	if err != nil {
		log.Fatalf("create google maps client: %s", err)
	}

	return MapsStore{
		client: c,
	}
}

func (s *MapsStore) GetLocationFromSearchText(search string) (*location.Location, error) {
	res := &location.Location{}

	r := &maps.TextSearchRequest{
		Query: search,
	}

	route, err := s.client.TextSearch(context.Background(), r)
	if err != nil {
		zap.S().Error(err)
		return res, err
	}

	if len(route.Results) == 0 {
		zap.S().Error(`Поиск не дал результатов`)

		return res, err
	}

	if len(route.Results) > 1 {
		zap.S().Error(`Слишком много результатов`)

		return res, err
	}

	res.Latitude = route.Results[0].Geometry.Location.Lat
	res.Longitude = route.Results[0].Geometry.Location.Lng

	return res, nil
}
