package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func checkPromoteParam(ctx context.Context, image, source, target string) (string, string) {
	if len(image) == 0 {
		logger.FatalfOnErr(ctx, errors.New("image is required"), "check image")
	}

	sourceValue := strings.TrimSpace(source)
	if len(sourceValue) == 0 {
		logger.FatalfOnErr(ctx, errors.New("source tag is required"), "check source")
	}

	targetValue := strings.TrimSpace(target)
	if len(targetValue) == 0 {
		logger.FatalfOnErr(ctx, errors.New("target tag is required"), "check target")
	}

	return sourceValue, targetValue
}

func promote(ctx context.Context, service RegistryService, image, source, target string) {
	manifest, err := service.GetManifest(ctx, image, source)
	if err != nil {
		logger.FatalfOnErr(ctx, fmt.Errorf("get manifest %s:%s: %w", image, source, err), "get manifest")
	}

	if err := service.PutManifest(ctx, image, target, manifest); err != nil {
		logger.FatalfOnErr(ctx, fmt.Errorf("put manifest %s:%s: %w", image, target, err), "put manifest")
	}
}
