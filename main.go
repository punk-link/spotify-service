package main

import (
	"main/models"
	"main/services"

	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"

	runtime "github.com/punk-link/streaming-platform-runtime"
	common "github.com/punk-link/streaming-platform-runtime/common"
	"github.com/punk-link/streaming-platform-runtime/startup"
)

func main() {
	logger := logger.New()
	environmentName := common.GetEnvironmentName()
	logger.LogInfo("%s is running as '%s'", SERVICE_NAME, environmentName)

	serviceOptions := runtime.NewServiceOptions(logger, environmentName, SERVICE_NAME)

	spotifySettingsValues, err := serviceOptions.Consul.Get("SpotifySettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
	}
	spotifySettings := spotifySettingsValues.(map[string]any)

	spotifyService := services.NewSpotifyService(logger, httpclient.DefaultConfig(logger), &models.SpotifyClientConfig{
		ClientId:     spotifySettings["ClientId"].(string),
		ClientSecret: spotifySettings["ClientSecret"].(string),
	})
	go startup.ProcessUrls(serviceOptions, spotifyService)

	startup.RunServer(serviceOptions)
}

const SERVICE_NAME = "spotify-service"
