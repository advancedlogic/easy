package interfaces

import (
	"time"
)

type Configuration interface {
	// Retrieve a configuration
	Open(...string) error
	// Save configuration
	Save(interface{}) error
	// Set based on key value
	Set(string, interface{}) error
	// Get based on key
	Get(string) (interface{}, error)

	GetIntOrDefault(string, int) int
	GetStringOrDefault(string, string) string
	GetBoolOrDefault(string, bool) bool
	GetFloat64OrDefault(string, float64) float64
	GetDurationOrDefault(string, time.Duration) time.Duration
	GetMapOfStringOrDefault(string, map[string]string) map[string]string
	GetArrayOfStringsOrDefault(string, []string) []string
}

type ConfigurationOption func(Configuration) error
