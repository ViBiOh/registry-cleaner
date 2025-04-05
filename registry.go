package main

import (
	"context"

	"github.com/ViBiOh/registry-cleaner/pkg/hub"
	"github.com/ViBiOh/registry-cleaner/pkg/registry"
	"github.com/docker/distribution"
)

const (
	dockerHub = "https://registry-1.docker.io/"
)

type RegistryService interface {
	Repositories(context.Context) ([]string, error)

	Tags(context.Context, string, func(string)) error
	Delete(context.Context, string, string) error

	GetManifest(ctx context.Context, repository, tag string) (distribution.Manifest, error)
	PutManifest(ctx context.Context, repository, tag string, manifest distribution.Manifest) error
}

func getRegistryService(ctx context.Context, registryURL, username, password, owner string) (RegistryService, error) {
	if registryURL == dockerHub {
		return hub.New(ctx, username, password, owner)
	}

	return registry.New(registryURL, username, password)
}
