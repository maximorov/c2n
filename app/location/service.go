package location

import (
	"fmt"
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
