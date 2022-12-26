package main

import (
	"errors"
	"main/models"
	"main/services"

	httpclient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"

	envManager "github.com/punk-link/environment-variable-manager"
	runtime "github.com/punk-link/streaming-platform-runtime"
	common "github.com/punk-link/streaming-platform-runtime/common"
	"github.com/punk-link/streaming-platform-runtime/startup"
)

func main() {
	logger := logger.New()
	envManager := envManager.New()

	environmentName := common.GetEnvironmentName(envManager)
	logger.LogInfo("%s is running as '%s'", SERVICE_NAME, environmentName)

	appSecrets := common.GetAppSecrets(envManager, logger, SECRET_ENGINE_NAME, SERVICE_NAME)
	serviceOptions := runtime.NewServiceOptions(logger, appSecrets, environmentName, SERVICE_NAME)

	spotifyClientSecret, isExist := appSecrets["client-secret"]
	if !isExist {
		err := errors.New("can't obtain host settings from Consul: '%s'")
		logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
	}

	spotifySettingsValues, err := serviceOptions.Consul.Get("SpotifySettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
	}
	spotifySettings := spotifySettingsValues.(map[string]any)

	spotifyService := services.NewSpotifyService(logger, httpclient.DefaultConfig(logger), &models.SpotifyClientConfig{
		ClientId:     spotifySettings["ClientId"].(string),
		ClientSecret: spotifyClientSecret.(string),
	})
	go startup.ProcessUrls(serviceOptions, spotifyService)

	startup.RunServer(serviceOptions)
}

const SECRET_ENGINE_NAME = "secrets"
const SERVICE_NAME = "spotify-service"
