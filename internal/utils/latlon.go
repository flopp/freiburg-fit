package utils

import (
	"math"

	"github.com/flopp/go-coordsparser"
)

type LatLon struct {
	Lat float64
	Lon float64
}

func ParseLatLon(input string) (LatLon, error) {
	lat, lon, err := coordsparser.Parse(input)
	if err != nil {
		// return invalid coordinates
		return LatLon{1000, 1000}, err
	}
	return LatLon{Lat: lat, Lon: lon}, nil
}

func (l LatLon) IsValid() bool {
	return -90.0 <= l.Lat && l.Lat <= 90.0 && -180.0 <= l.Lon && l.Lon <= 180.0
}

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func Distance(aa LatLon, bb LatLon) float64 {

	const earthRadiusKM float64 = 6371.0

	lat1 := deg2rad(aa.Lat)
	lon1 := deg2rad(aa.Lon)
	lat2 := deg2rad(bb.Lat)
	lon2 := deg2rad(bb.Lon)

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	distance := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)) * earthRadiusKM

	return distance
}
