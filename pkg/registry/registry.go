package registry

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/heroku/docker-registry-client/registry"
)

// App of package
type App struct {
	client *registry.Registry
}

// New creates new App from Config
func New(url, username, password string) (App, error) {
	client, err := registry.New(url, username, password)
	if err != nil {
		return App{}, err
	}

	client.Logf = func(format string, args ...interface{}) {
		slog.Info(fmt.Sprintf(format, args))
	}

	return App{
		client: client,
	}, nil
}

// Repositories list repositories
func (a App) Repositories(_ context.Context) ([]string, error) {
	return a.client.Repositories()
}

// Tags list tags for a given image
func (a App) Tags(ctx context.Context, image string, handler func(string)) error {
	tags, err := a.client.Tags(image)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		handler(tag)
	}

	return nil
}

// Delete a tag
func (a App) Delete(_ context.Context, image, tag string) error {
	digest, err := a.client.ManifestDigest(image, tag)
	if err != nil {
		return fmt.Errorf("get manifest: %w", err)
	}

	if err = a.client.DeleteManifest(image, digest); err != nil {
		return fmt.Errorf("delete manifest: %w", err)
	}

	return nil
}
