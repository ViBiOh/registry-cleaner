package main

import (
	"context"
	"log/slog"
	"os"
	"regexp"
	"runtime"

	"github.com/ViBiOh/httputils/v4/pkg/concurrent"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func deleteTags(ctx context.Context, service RegistryService, image string, matcher *regexp.Regexp, tagIndex int, last, invert, delete bool) {
	lastTags := make(map[string]string)
	limiter := concurrent.NewLimiter(runtime.NumCPU())

	err := service.Tags(ctx, image, func(tag string) {
		tagBucket, found := getTagAndMatchTag(tag, matcher, tagIndex)
		if !found {
			return
		}

		if last {
			if lastHandler(ctx, service, invert, delete, image, tag, tagBucket, lastTags) {
				return
			}
		}

		limiter.Go(func() {
			tagHandler(ctx, service, delete, image, tag)
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
