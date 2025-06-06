package utils

import (
	"log"

	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	Sugar = logger.Sugar()
	defer logger.Sync()
}
