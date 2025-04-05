package main

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func main() {
	config := newConfig()

	ctx := context.Background()
	logger.Init(ctx, config.logger)

	registryURL, imageName, matcher, tagIndex := checkParam(ctx, *config.url, *config.image, *config.grep, *config.list)

	service, err := getRegistryService(ctx, registryURL, *config.username, *config.password, *config.owner)
	logger.FatalfOnErr(ctx, err, "create registry client")

	if *config.list {
		listRepositories(ctx, service, imageName)
		return
	}

	deleteTags(ctx, service, imageName, matcher, tagIndex, *config.last, *config.invert, *config.delete)
}

func checkParam(ctx context.Context, url, image, grep string, list bool) (string, string, *regexp.Regexp, int) {
	registryURL := strings.TrimSpace(url)
	if len(registryURL) == 0 {
		logger.FatalfOnErr(ctx, errors.New("url is required"), "check url")
	}

	imageName := strings.ToLower(strings.TrimSpace(image))

	if list {
		return registryURL, imageName, nil, 0
	}

	if len(imageName) == 0 {
		logger.FatalfOnErr(ctx, errors.New("image is required"), "check image")
	}

	grepValue := strings.TrimSpace(grep)
	if len(grepValue) == 0 {
		logger.FatalfOnErr(ctx, errors.New("grep pattern is required"), "check grep")
	}

	matcher, err := regexp.Compile(grepValue)
	logger.FatalfOnErr(ctx, err, "compile grep regexp")

	tagIndex := -1

	for i, group := range matcher.SubexpNames() {
		if group == "tagBucket" {
			tagIndex = i
			break
		}
	}

	return registryURL, imageName, matcher, tagIndex
}
