package main

import (
	"avito_app/internal/models"
	"time"
)

type cachedData map[pair]cacheValue

type cacheValue struct {
	banner    *models.Banner
	timestamp time.Time
}

type pair struct {
	tagID     int
	featureID int
}
