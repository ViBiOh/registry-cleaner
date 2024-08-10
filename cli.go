package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/concurrent"
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
		listRepositories(ctx, service)
		return
	}

	lastTags := make(map[string]string)
	limiter := concurrent.NewLimiter(runtime.NumCPU())

	err = service.Tags(ctx, imageName, func(tag string) {
		tagBucket, found := getTagAndMatchTag(tag, matcher, tagIndex)
		if !found {
			return
		}

		if *config.last {
			if lastHandler(ctx, service, *config.invert, *config.delete, imageName, tag, tagBucket, lastTags) {
				return
			}
		}

		limiter.Go(func() {
			tagHandler(ctx, service, *config.delete, imageName, tag)
		})
	})

	logger.FatalfOnErr(ctx, err, "list tags")

	limiter.Wait()
}

func getTagAndMatchTag(tag string, matcher *regexp.Regexp, tagIndex int) (string, bool) {
	matches := matcher.FindStringSubmatch(tag)
	if len(matches) == 0 {
		return "", false
	}

	if tagIndex > 0 {
		if tagIndex < len(matches) {
			return matches[tagIndex], true
		}

		return "", false
	}

	return "", true
}

func lastHandler(ctx context.Context, service RegistryService, invert, delete bool, image, tag, tagBucket string, lastTags map[string]string) bool {
	if len(lastTags[tagBucket]) == 0 {
		lastTags[tagBucket] = tag

		return true
	}

	lastTag := lastTags[tagBucket]

	if (invert && tag < lastTag) || tag > lastTag {
		tagHandler(ctx, service, delete, image, lastTag)

		lastTags[tagBucket] = tag

		return true
	}

	return false
}

func tagHandler(ctx context.Context, service RegistryService, delete bool, image, tag string) {
	if !delete {
		slog.LogAttrs(ctx, slog.LevelWarn, "eligible to deletion", slog.String("image", image), slog.String("tag", tag))
		return
	}

	if err := service.Delete(ctx, image, tag); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "delete", slog.String("image", image), slog.String("tag", tag), slog.Any("error", err))
		os.Exit(1)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "deleted", slog.String("image", image), slog.String("tag", tag))
}

func checkParam(ctx context.Context, url, image, grep string, list bool) (string, string, *regexp.Regexp, int) {
	registryURL := strings.TrimSpace(url)
	if len(registryURL) == 0 {
		logger.FatalfOnErr(ctx, errors.New("url is required"), "check url")
	}

	if list {
		return registryURL, "", nil, 0
	}

	imageName := strings.ToLower(strings.TrimSpace(image))
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
