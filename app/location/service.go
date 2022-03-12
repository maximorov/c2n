package location

import (
	"fmt"
	"math"
)

type (
	Store interface {
		GetOne() (Location, error)
		// CountDistance returns distance in meters between two locations.
		CountDistance(location1, location2 Location) uint
		// FoundTasksInRadius find locations within a radius(in meters)
		//TODO: change locations to tasks
		FoundTasksInRadius(location Location, distance float64) ([]Location, error)
	}

	storeImpl struct {
	}
)

const PI float64 = 3.141592653589793

func (s *storeImpl) GetOne(longitude, latitude float64) (Location, error) {
	if longitude == 0 || latitude == 0 {
		return Location{}, fmt.Errorf("not correct cordinates")
	}

	return Location{
		Longitude: longitude,
		Latitude:  latitude,
	}, nil
}

func (s *storeImpl) CountDistance(loc1, loc2 Location) float64 {
	radlat1 := float64(PI * loc1.Latitude / 180)
	radlat2 := float64(PI * loc2.Latitude / 180)

	theta := float64(loc1.Longitude - loc2.Longitude)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515 * 1609.344

	return dist
}

func (s *storeImpl) FoundTasksInRadius(location Location, distance float64) ([]Location, error) {
	var result []Location
	//TODO: change to get tasks
	tasks := []Location{{Longitude: 50.385896, Latitude: 30.464168}}

	for _, task := range tasks {
		dist := s.CountDistance(task, location)
		if dist < distance {
			result = append(result, task)
		}

	}

	return result, nil
}
