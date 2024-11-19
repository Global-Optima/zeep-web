package utils

func IsValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func IsValidLongitude(lon float64) bool {
	return lon >= -180 && lon <= 180
}
