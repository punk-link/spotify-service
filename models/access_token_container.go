package models

import "time"

type SpotifyAccessTokenContainer struct {
	Token   string
	Expired time.Time
}
