package models

type ArtistReleasesContainer struct {
	Items []Release `json:"items"`
	Next  string    `json:"next"`
}
