package services

import (
	"encoding/base64"
	"main/models"
	"net/http"
	"net/url"
	"strings"
	"time"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
)

func getAccessToken(logger logger.Logger, httpClient httpClient.HttpClient[models.SpotifyAccessToken], spotifyConfig *models.SpotifyClientConfig) (string, error) {
	if len(_tokenContainer.Token) != 0 && time.Now().UTC().Before(_tokenContainer.Expired) {
		return _tokenContainer.Token, nil
	}

	request, err := getAccessTokenRequest(logger, spotifyConfig)
	if err != nil {
		logger.LogError(err, err.Error())
		return "", err
	}

	accessToken, err := httpClient.MakeRequest(request)
	if err != nil {
		logger.LogError(err, err.Error())
		return "", err
	}

	_tokenContainer = models.SpotifyAccessTokenContainer{
		Expired: time.Now().Add(time.Second*time.Duration(accessToken.ExpiresIn) - ACCESS_TOKEN_SAFITY_THRESHOLD).UTC(),
		Token:   accessToken.Token,
	}

	return _tokenContainer.Token, nil
}

func getAccessTokenRequest(logger logger.Logger, config *models.SpotifyClientConfig) (*http.Request, error) {
	payload := url.Values{}
	payload.Add("grant_type", "client_credentials")

	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(payload.Encode()))
	if err != nil {
		logger.LogError(err, err.Error())
		return nil, err
	}

	credentials := "Basic " + base64.StdEncoding.EncodeToString([]byte(config.ClientId+":"+config.ClientSecret))

	request.Header.Add("Authorization", credentials)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

var _tokenContainer models.SpotifyAccessTokenContainer

const ACCESS_TOKEN_SAFITY_THRESHOLD = time.Second * 5
