package main

import (
	"context"
	"errors"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func main() {
	config := newConfig()

	ctx := context.Background()
	logger.Init(ctx, config.logger)

	registryURL, image := checkParam(ctx, *config.url, *config.image)

	service, err := getRegistryService(ctx, registryURL, *config.username, *config.password, *config.owner)
	logger.FatalfOnErr(ctx, err, "create registry client")

	switch config.command {
	case "list":
		listRepositories(ctx, service, image)

	case "promote":
		source, target := checkPromoteParam(ctx, image, *config.source, *config.target)
		promote(ctx, service, image, source, target)

	case "delete":
		matcher, tagIndex := checkDeleteParam(ctx, image, *config.grep)
		deleteTags(ctx, service, image, matcher, tagIndex, *config.last, *config.invert, *config.dryRun)
	}
}

func checkParam(ctx context.Context, url, image string) (string, string) {
	registryURL := strings.TrimSpace(url)
	if len(registryURL) == 0 {
		logger.FatalfOnErr(ctx, errors.New("url is required"), "check url")
	}

	imageName := strings.ToLower(strings.TrimSpace(image))

	return registryURL, imageName
}
