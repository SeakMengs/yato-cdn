package util

import (
	"fmt"
	"math"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// calculates distance between two lat/lon points
// credit: https://gist.github.com/hotdang-ca/6c1ee75c48e515aec5bc6db6e3265e49
// :::    optional: unit = the unit you desire for results                     :::
// :::           where: 'M' is statute miles (default, or omitted)             :::
// :::                  'K' is kilometers                                      :::
// :::                  'N' is nautical miles                                  :::
func Distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
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
