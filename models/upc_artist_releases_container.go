package models

type UpcArtistReleasesContainer struct {
	Albums ArtistReleasesContainer `json:"albums"`
}
