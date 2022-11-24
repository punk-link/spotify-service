package services

import (
	"fmt"
	"main/models"
	"net/http"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
)

func makeBatchRequestWithSync[T any](logger logger.Logger, httpClientConfig *httpClient.HttpClientConfig, spotifyConfig *models.SpotifyClientConfig, upcContainers []platformContracts.UpcContainer) []httpClient.SyncedResult[T] {
	syncedHttpRequests := make([]httpClient.SyncedRequest, len(upcContainers))
	for i, upcContainer := range upcContainers {
		request, err := getUpcRequest(logger, httpClientConfig, spotifyConfig, upcContainer.Upc)
		if err != nil {
			logger.LogWarn("can't build an http request: %s", err.Error())
			continue
		}

		syncedHttpRequests[i] = httpClient.SyncedRequest{
			HttpRequest: request,
			SyncKey:     upcContainer.Upc,
		}
	}

	return httpClient.MakeBatchRequestWithSyncKeys[T](httpClient.DefaultConfig(logger), syncedHttpRequests)
}

func getUpcRequest(logger logger.Logger, httpClientConfig *httpClient.HttpClientConfig, spotifyConfig *models.SpotifyClientConfig, url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/search?type=album&q=upc:%s", url), nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken(logger, httpClientConfig, spotifyConfig)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
