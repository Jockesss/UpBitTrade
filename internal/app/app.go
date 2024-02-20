package app

import (
	"fmt"
	"upbit/internal/config"
	"upbit/pkg/log"
)

func Run() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to load config: %v", err))
	}

	log.Logger.Info(fmt.Sprintf("Config FILE --> ", cfg))

}
