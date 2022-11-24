package services

import (
	"main/models"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
)

type SpotifyService struct {
	spotifyConfig    *models.SpotifyClientConfig
	httpClientConfig *httpClient.HttpClientConfig
	logger           logger.Logger
}

func NewSpotifyService(logger logger.Logger, httpClientConfig *httpClient.HttpClientConfig, spotifyConfig *models.SpotifyClientConfig) *SpotifyService {
	return &SpotifyService{
		spotifyConfig:    spotifyConfig,
		httpClientConfig: httpClientConfig,
		logger:           logger,
	}
}

func (t *SpotifyService) GetBatchSize() int {
	return 40
}

func (t *SpotifyService) GetPlatformName() string {
	return platformContracts.Spotify
}

func (t *SpotifyService) GetReleaseUrlsByUpc(upcContainers []platformContracts.UpcContainer) []platformContracts.UrlResultContainer {
	syncedReleaseContainers := makeBatchRequestWithSync[models.UpcArtistReleasesContainer](t.logger, t.httpClientConfig, t.spotifyConfig, upcContainers)

	upcMap := t.getUpcMap(upcContainers)
	results := make([]platformContracts.UrlResultContainer, 0)
	for _, syncedContainer := range syncedReleaseContainers {
		container := syncedContainer.Result
		if len(container.Albums.Items) == 0 {
			continue
		}

		id := upcMap[syncedContainer.SyncKey]
		results = append(results, platformContracts.UrlResultContainer{
			Id:           id,
			PlatformName: t.GetPlatformName(),
			Upc:          syncedContainer.SyncKey,
			Url:          container.Albums.Items[0].ExternalUrls.Spotify,
		})
	}

	return results
}

func (t *SpotifyService) getUpcMap(upcContainers []platformContracts.UpcContainer) map[string]int {
	results := make(map[string]int, len(upcContainers))
	for _, container := range upcContainers {
		results[container.Upc] = container.Id
	}

	return results
}
