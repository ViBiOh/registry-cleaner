package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func listRepositories(ctx context.Context, service RegistryService, image string) {
	if len(image) > 0 {
		if err := service.Tags(ctx, image, func(tag string) { fmt.Println(tag) }); err != nil {
			logger.FatalfOnErr(ctx, err, "list tags")
		}

		return
	}

	repositories, err := service.Repositories(ctx)
	logger.FatalfOnErr(ctx, err, "list repositories")

	fmt.Printf("%s", strings.Join(repositories, "\n"))
}
