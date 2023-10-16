package main

import "os"

func Getenv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetBoolEnv(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}
	return fallback
}

func msToSeconds(ms float64) float64 {
	return ms / 1000
}

func percentageToRatio(percentage float64) float64 {
	return percentage / 100
}
