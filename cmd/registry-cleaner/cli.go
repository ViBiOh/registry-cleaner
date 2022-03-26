package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ViBiOh/flags"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/registry-cleaner/pkg/hub"
	"github.com/ViBiOh/registry-cleaner/pkg/registry"
)

const (
	dockerHub = "https://registry-1.docker.io/"
)

// RegistryService definition
type RegistryService interface {
	Tags(context.Context, string, func(string)) error
	Delete(context.Context, string, string) error
}

func main() {
	fs := flag.NewFlagSet("registry-cleaner", flag.ExitOnError)

	loggerConfig := logger.Flags(fs, "logger")

	url := flags.String(fs, "", "registry", "URL", "Registry URL", dockerHub, nil)
	username := flags.String(fs, "", "registry", "Username", "Registry username", "", nil)
	password := flags.String(fs, "", "registry", "Password", "Registry password", "", nil)
	image := flags.String(fs, "", "registry", "Image", "Image name", "", nil)
	grep := flags.String(fs, "", "registry", "Grep", "Matching tags regexp", "", nil)
	dryRun := flags.Bool(fs, "", "cleaner", "DryRun", "Dry run mode", false, nil)

	logger.Fatal(fs.Parse(os.Args[1:]))

	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	registryURL := strings.TrimSpace(*url)
	if len(registryURL) == 0 {
		logger.Fatal(errors.New("-url flag is required"))
	}

	imageName := strings.TrimSpace(*image)
	if len(imageName) == 0 {
		logger.Fatal(errors.New("-image flag is required"))
	}

	grepValue := strings.TrimSpace(*grep)
	if len(grepValue) == 0 {
		logger.Fatal(errors.New("-grep flag is required"))
	}

	matcher, err := regexp.Compile(grepValue)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to compile grep regexp: %s", err))
	}

	var service RegistryService
	if registryURL == dockerHub {
		service, err = hub.New(context.Background(), *username, *password)
	} else {
		service, err = registry.New(registryURL, *username, *password)
	}

	if err != nil {
		logger.Fatal(fmt.Errorf("unable to create registry client: %s", err))
	}

	err = service.Tags(context.Background(), imageName, func(tag string) {
		if !matcher.MatchString(tag) {
			return
		}

		if *dryRun {
			logger.Warn("%s:%s is matching", imageName, tag)
			return
		}

		err = service.Delete(context.Background(), imageName, tag)
		if err != nil {
			logger.Fatal(fmt.Errorf("unable to delete `%s:%s`: %s", imageName, tag, err))
		}
		logger.Info("%s:%s deleted!", imageName, tag)
	})

	if err != nil {
		logger.Fatal(fmt.Errorf("unable to list tags: %s", err))
	}
}
