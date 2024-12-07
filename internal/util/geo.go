package util

import (
	"fmt"
	"math"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// calculates distance between two lat/lon points
func Haversine(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const R = 6371 // Radius of the Earth in km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Region    string  `json:"region"`
}

func FindDomainIp(domain string) string {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return ""
	}
	if len(ips) == 0 {
		return ""
	}

	// use the first ip address for geolocation
	ip := ips[0].String()
	fmt.Printf("ip address of %s: %s\n", domain, ip)

	return ip
}

//	func FindGeoLocation(ip string) (*GeoLocation, error) {
//		geoIPAPI := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
//		resp, err := http.Get(geoIPAPI)
//		if err != nil {
//			return nil, err
//		}
//		defer resp.Body.Close()
//
//		var geo *GeoLocation
//		if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
//			return nil, err
//		}
//
//		return geo, nil
//	}

func FindGeoLocation(ip string) (*GeoLocation, error) {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, err
	}

	record, err := db.City(parsedIP)
	if err != nil {
		return nil, err
	}

	// Extract location details
	geo := &GeoLocation{
		Latitude:  record.Location.Latitude,
		Longitude: record.Location.Longitude,
	}

	geo.Region = record.Country.Names["en"]
	// if len(record.Subdivisions) > 0 {
	// 	geo.Region = record.Subdivisions[0].Names["en"]
	// } else {
	// 	// Use country as a fallback for region
	// }

	return geo, nil
}
