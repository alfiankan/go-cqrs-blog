package config

import "go.uber.org/zap"

// InitLog init zap log by config app_env
// production [info, warn, err]
// development [debug, info, warn, err]
func InitLog(config ApplicationConfig) (logger *zap.Logger, err error) {

	defer logger.Sync()
	if config.AppEnv == "development" {
		logger, err = zap.NewDevelopment()
		return
	}

	logger, err = zap.NewProduction()
	return
}
