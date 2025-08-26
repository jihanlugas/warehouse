package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jihanlugas/warehouse/config"
)

func GetUniqueID() string {
	uuid := uuid.New().String()
	return uuid
}

// GetRandomNumber returns a random integer between min and max (inclusive)
func GetRandomNumber(min, max int) int {
	if min > max {
		min, max = max, min // Swap if min is greater than max
	}

	// Create a new random source and generator
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	return r.Intn(max-min+1) + min
}

func GetPhotoUrlById(photoID string) string {
	return fmt.Sprintf("%s/photo/%s", config.Server.BaseUrl, photoID)
}
