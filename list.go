package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/logger"
)

func listRepositories(ctx context.Context, service RegistryService) {
	repositories, err := service.Repositories(ctx)
	logger.FatalfOnErr(ctx, err, "list repositories")

	fmt.Printf("%s", strings.Join(repositories, "\n"))
}
