package log

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var err error
	// Init logger
	//Logger, err = zap.NewProduction() // For PROD
	Logger, err = zap.NewDevelopment() // For DEV

	if err != nil {
		panic(err)
	}
}
