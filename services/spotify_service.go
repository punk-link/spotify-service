package services

import (
	"fmt"
	"main/models"
	"net/http"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
)

type SpotifyService struct {
	spotifyConfig   *models.SpotifyClientConfig
	tokenHttpClient httpClient.HttpClient[models.SpotifyAccessToken]
	upcHttpClient   httpClient.HttpClient[models.UpcArtistReleasesContainer]
	logger          logger.Logger
}

func NewSpotifyService(logger logger.Logger, httpClientConfig *httpClient.HttpClientConfig, spotifyConfig *models.SpotifyClientConfig) *SpotifyService {
	tokenHttpClient := httpClient.New[models.SpotifyAccessToken](httpClientConfig)
	upcHttpClient := httpClient.New[models.UpcArtistReleasesContainer](httpClientConfig)

	return &SpotifyService{
		spotifyConfig:   spotifyConfig,
		tokenHttpClient: tokenHttpClient,
		upcHttpClient:   upcHttpClient,
		logger:          logger,
	}
}

func (t *SpotifyService) GetBatchSize() int {
	return 40
}

func (t *SpotifyService) GetPlatformName() string {
	return platformContracts.Spotify
}

func (t *SpotifyService) GetReleaseUrlsByUpc(upcContainers []platformContracts.UpcContainer) []platformContracts.UrlResultContainer {
	syncedReleaseContainers := t.makeBatchRequestWithSync(upcContainers)

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

func (t *SpotifyService) makeBatchRequestWithSync(upcContainers []platformContracts.UpcContainer) []httpClient.SyncedResult[models.UpcArtistReleasesContainer] {
	syncedHttpRequests := make([]httpClient.SyncedRequest, len(upcContainers))
	for i, upcContainer := range upcContainers {
		request, err := getUpcRequest(t.logger, t.tokenHttpClient, t.spotifyConfig, upcContainer.Upc)
		if err != nil {
			t.logger.LogWarn("can't build an http request: %s", err.Error())
			continue
		}

		syncedHttpRequests[i] = httpClient.SyncedRequest{
			HttpRequest: request,
			SyncKey:     upcContainer.Upc,
		}
	}

	return t.upcHttpClient.MakeBatchRequestWithSync(syncedHttpRequests)
}

func getUpcRequest(logger logger.Logger, httpClient httpClient.HttpClient[models.SpotifyAccessToken], spotifyConfig *models.SpotifyClientConfig, url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/search?type=album&q=upc:%s", url), nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken(logger, httpClient, spotifyConfig)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
