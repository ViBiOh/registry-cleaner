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

func checkDeleteParam(ctx context.Context, image, grep string) (*regexp.Regexp, int) {
	if len(image) == 0 {
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

	return matcher, tagIndex
}

func deleteTags(ctx context.Context, service RegistryService, image string, matcher *regexp.Regexp, tagIndex int, last, invert, dryRun bool) {
	lastTags := make(map[string]string)
	limiter := concurrent.NewLimiter(runtime.NumCPU())

	err := service.Tags(ctx, image, func(tag string) {
		tagBucket, found := getTagAndMatchTag(tag, matcher, tagIndex)
		if !found {
			return
		}

		if last {
			if lastHandler(ctx, service, invert, dryRun, image, tag, tagBucket, lastTags) {
				return
			}
		}

		limiter.Go(func() {
			tagHandler(ctx, service, dryRun, image, tag)
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

func lastHandler(ctx context.Context, service RegistryService, invert, dryRun bool, image, tag, tagBucket string, lastTags map[string]string) bool {
	if len(lastTags[tagBucket]) == 0 {
		lastTags[tagBucket] = tag

		return true
	}

	lastTag := lastTags[tagBucket]

	if (invert && tag < lastTag) || tag > lastTag {
		tagHandler(ctx, service, dryRun, image, lastTag)

		lastTags[tagBucket] = tag

		return true
	}

	return false
}

func tagHandler(ctx context.Context, service RegistryService, dryRun bool, image, tag string) {
	if dryRun {
		slog.LogAttrs(ctx, slog.LevelWarn, "eligible to deletion", slog.String("image", image), slog.String("tag", tag))
		return
	}

	if err := service.Delete(ctx, image, tag); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "delete", slog.String("image", image), slog.String("tag", tag), slog.Any("error", err))
		os.Exit(1)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "deleted", slog.String("image", image), slog.String("tag", tag))
}
